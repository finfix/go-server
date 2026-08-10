package service

import (
	"context"

	"github.com/google/uuid"
)

// GetAccessedAccountGroups возвращает доступы пользователей к группам счетов
func (s *UserToAccountGroupService) GetAccessedAccountGroups(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "GetAccessedAccountGroups")
	defer span.End()

	return s.userToAccountGroupRepository.GetAccessedAccountGroups(ctx, userID)
}
