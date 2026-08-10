package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"

	"server/internal/modules/userToAccountGroup/repository"
)

var tracer = otel.Tracer("/server/internal/modules/userToAccountGroup/service")

// UserToAccountGroupService - сервис связей пользователей с группами счетов
type UserToAccountGroupService struct {
	userToAccountGroupRepository UserToAccountGroupRepository
}

var _ UserToAccountGroupRepository = new(repository.UserToAccountGroupRepository)

// UserToAccountGroupRepository - интерфейс репозитория связей пользователей с группами счетов
type UserToAccountGroupRepository interface {
	GetAccessedAccountGroups(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

// NewUserToAccountGroupService создает новый сервис связей пользователей с группами счетов
func NewUserToAccountGroupService(
	userToAccountGroupRepository UserToAccountGroupRepository,
) *UserToAccountGroupService {
	return &UserToAccountGroupService{
		userToAccountGroupRepository: userToAccountGroupRepository,
	}
}
