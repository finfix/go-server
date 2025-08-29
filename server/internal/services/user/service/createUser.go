package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/services/user/model"
)

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(ctx context.Context, user model.CreateReq) (id uuid.UUID, err error) {
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()

	return s.userRepository.CreateUser(ctx, user)
}
