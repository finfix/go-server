package service

import (
	"context"

	"pkg/slices"
	"server/internal/utils/errors"

	userModel "server/internal/modules/user/model"

	"github.com/google/uuid"
)

func (s *SettingsService) checkAdmin(ctx context.Context, userID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "checkAdmin")
	defer span.End()

	// Получаем пользователя по ID
	user, err := slices.FirstWithError(
		s.userService.GetUsers(ctx, userModel.GetUsersReq{ //nolint:exhaustruct
			IDs: []uuid.UUID{userID},
		}),
	)
	if err != nil {
		return err
	}

	// Проверяем, является ли пользователь администратором
	if !user.IsAdmin {
		return errors.Forbidden.New("Access denied").WithContextParams(ctx)
	}

	return nil
}
