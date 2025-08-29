package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/modules/user/model"
)

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(ctx context.Context, user model.CreateReq) (id uuid.UUID, err error) {
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()

	user.ID = uuid.New()

	return user.ID, s.userRepository.CreateUser(ctx, user)
}
