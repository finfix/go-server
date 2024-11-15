package service

import (
	"context"

	"pkg/errors"
	"pkg/slices"

	userModel "server/internal/services/user/model"
)

func (s *SettingsService) checkAdmin(ctx context.Context, userID uint32) error {
	ctx, span := tracer.Start(ctx, "checkAdmin")
	defer span.End()

	// Получаем пользователя по ID
	user, err := slices.FirstWithError(
		s.userService.GetUsers(ctx, userModel.GetUsersReq{ //nolint:exhaustruct
			IDs: []uint32{userID},
		}),
	)
	if err != nil {
		return err
	}

	// Проверяем, является ли пользователь администратором
	if !user.IsAdmin {
		return errors.Forbidden.New("Access denied")
	}

	return nil
}
