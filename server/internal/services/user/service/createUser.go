package service

import (
	"context"

	"server/internal/services/user/model"
)

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(ctx context.Context, user model.CreateReq) (id uint32, err error) {
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()

	return s.userRepository.CreateUser(ctx, user)
}
