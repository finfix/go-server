package service

import (
	"context"

	"pkg/errors"
	"pkg/passwordManager"

	"server/internal/services/user/model"
	userRepoModel "server/internal/services/user/repository/model"
)

// UpdateUser обновляет настройки пользователя
func (s *UserService) UpdateUser(ctx context.Context, req model.UpdateUserReq) error {
	ctx, span := tracer.Start(ctx, "UpdateUser")
	defer span.End()

	return s.generalRepository.WithinTransaction(ctx, func(ctx context.Context) error {

		// Если обновляется токен уведомлений, обновляем его в таблице девайсов
		if req.NotificationToken != nil {
			if err := s.userRepository.UpdateDevice(ctx, userRepoModel.UpdateDeviceReq{
				UserID:            req.Necessary.UserID,
				DeviceID:          req.Necessary.DeviceID,
				RefreshToken:      nil,
				NotificationToken: req.NotificationToken,
				ApplicationInformation: userRepoModel.UpdateApplicationInformationReq{
					BundleID: nil,
					Version:  nil,
					Build:    nil,
				},
				DeviceInformation: userRepoModel.UpdateDeviceInformationReq{
					VersionOS: nil,
					IPAddress: nil,
					UserAgent: nil,
				},
			}); err != nil {
				return err
			}
		}

		repoReq := req.ConvertToRepoModel()

		// Если обновляется пароль
		if req.Password != nil {

			if req.OldPassword != nil {
				return errors.BadRequest.New("OldPassword must be filled")
			}

			// Получаем актуальный пароль пользователя
			users, err := s.userRepository.GetUsers(ctx, model.GetUsersReq{ //nolint:exhaustruct
				IDs: []uint32{req.Necessary.UserID},
			})
			if err != nil {
				return err
			}
			if len(users) == 0 {
				return errors.NotFound.New("User not found")
			}
			user := users[0]

			// Сравниваем пришедший пароль и хэш пароля из базы данных
			if err = passwordManager.CompareHashAndPassword(user.PasswordHash, []byte(*req.OldPassword), user.PasswordSalt, s.generalSalt); err != nil {
				return err
			}

			// Получаем хэш и соль нового пароля
			passwordHash, passwordSalt, err := passwordManager.CreateNewPassword([]byte(*req.Password), s.generalSalt)
			if err != nil {
				return err
			}

			repoReq.PasswordHash = &passwordHash
			repoReq.PasswordSalt = &passwordSalt
		}

		if err := s.userRepository.UpdateUser(ctx, repoReq); err != nil {
			return err
		}

		return nil
	})
}
