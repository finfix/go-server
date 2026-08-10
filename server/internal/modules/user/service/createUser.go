package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/user/model"
)

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(ctx context.Context, user model.CreateReq) (id uuid.UUID, err error) {
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()

	user.ID = uuid.New()

	err = s.generalRepository.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Создаем пользователя
		if err := s.userRepository.CreateUser(ctxTx, user); err != nil {
			return err
		}

		// Фиксируем создание пользователя в аудит-логе (без хэша и соли пароля)
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.User,
			Method:   auditLogMethod.Create,
			EntityID: user.ID.String(),
			Before:   nil,
			After: struct {
				ID              uuid.UUID
				Name            string
				Email           string
				DefaultCurrency string
				IsAdmin         bool
			}{user.ID, user.Name, user.Email, user.DefaultCurrency, user.IsAdmin},
			UserID: user.ID,
		})
	})
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}
