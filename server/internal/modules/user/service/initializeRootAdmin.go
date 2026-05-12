package service

import (
	"context"
	"time"

	"pkg/log"
	"pkg/slices"

	userModel "server/internal/modules/user/model"

	"github.com/google/uuid"
)

func (s *UserService) initializeRootAdmin(ctx context.Context) error {
	// Ищем пользователя с нулевым UUID
	_, err := slices.FirstWithError(s.userRepository.GetUsers(ctx, userModel.GetUsersReq{ //nolint:exhaustruct
		IDs: []uuid.UUID{uuid.Nil},
	}))
	if err == nil {
		log.WithParams("userID", uuid.Nil).Info("Root admin already exists")
		return nil
	}

	// Создаём админа с нулевым UUID, без пароля и логина
	err = s.userRepository.CreateUser(ctx, userModel.CreateReq{
		ID:              uuid.Nil,
		Name:            "Admin",
		Email:           "",
		PasswordHash:    []byte{},
		PasswordSalt:    []byte{},
		TimeCreate:      time.Now(),
		DefaultCurrency: "RUB",
		IsAdmin:         true,
	})
	if err != nil {
		return err
	}

	log.WithParams("userID", uuid.Nil).Info("Root admin created")

	return nil
}
