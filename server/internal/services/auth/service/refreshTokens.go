package service

import (
	"context"

	"github.com/google/uuid"
	"pkg/jwtManager"
	"server/internal/utils/auth"
	"server/internal/utils/errors"

	"server/internal/services/auth/model"
	"server/internal/services/auth/service/utils"
	userRepoModel "server/internal/services/user/repository/model"
)

// RefreshTokens обновляет токены доступа в базе данных
func (s *AuthService) RefreshTokens(ctx context.Context, req model.RefreshTokensReq) (newTokens model.RefreshTokensRes, err error) {
	ctx, span := tracer.Start(ctx, "RefreshTokens")
	defer span.End()

	// Получаем девайс по идентификатору пользователя и девайса
	devices, err := s.userRepository.GetDevices(ctx, userRepoModel.GetDevicesReq{ // nolint:exhaustruct
		DeviceIDs: []string{req.Necessary.DeviceID},
		UserIDs:   []uuid.UUID{req.Necessary.UserID},
	})
	if err != nil {
		return newTokens, err
	}
	if len(devices) == 0 {
		return newTokens, errors.Unauthorized.New("Device not found").WithContextParams(ctx).
			WithCustomHumanText("Девайс не найден")
	}
	device := devices[0]

	// Сравниваем токен из базы данных с переданным пользователем токеном
	if req.Token != device.RefreshToken {
		return newTokens, errors.Forbidden.New("Auth is incorrect").
			WithContextParams(ctx)
	}

	// Парсим токен
	claims, err := jwtManager.ParseToken[auth.Claims](device.RefreshToken, jwtManager.RefreshToken)
	if err != nil {
		return newTokens, err
	}

	// Дополнительно проверяем идентификаторы
	if claims.UserID != req.Necessary.UserID {
		return newTokens, errors.Forbidden.New("UserID not matched").
			WithContextParams(ctx)
	}
	if claims.DeviceID != req.Necessary.DeviceID {
		return newTokens, errors.Forbidden.New("DeviceID not matched").
			WithContextParams(ctx)
	}

	// Создаем новую пару токенов
	newTokens.Tokens, err = utils.CreatePairTokens(req.Necessary.UserID, req.Necessary.DeviceID)
	if err != nil {
		return newTokens, err
	}

	// Обновляем данные у девайса
	if err = s.userRepository.UpdateDevice(ctx, userRepoModel.UpdateDeviceReq{
		UserID:       req.Necessary.UserID,
		DeviceID:     req.Necessary.DeviceID,
		RefreshToken: &newTokens.Tokens.RefreshToken,
		DeviceInformation: userRepoModel.UpdateDeviceInformationReq{
			VersionOS: &req.Device.VersionOS,
			UserAgent: &req.Device.UserAgent,
			IPAddress: &req.Device.IPAddress,
		},
		ApplicationInformation: userRepoModel.UpdateApplicationInformationReq{
			BundleID: &req.Application.BundleID,
			Version:  &req.Application.Version,
			Build:    &req.Application.Build,
		},
		NotificationToken: nil,
	}); err != nil {
		return newTokens, err
	}

	// Возвращаем пару токенов клиенту
	return newTokens, nil
}
