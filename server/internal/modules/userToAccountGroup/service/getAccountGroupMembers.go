package service

import (
	"context"

	"github.com/google/uuid"
)

// GetAccountGroupMembers возвращает всех пользователей, имеющих доступ к группе счетов.
func (s *UserToAccountGroupService) GetAccountGroupMembers(ctx context.Context, accountGroupID uuid.UUID) ([]uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "GetAccountGroupMembers")
	defer span.End()

	return s.userToAccountGroupRepository.GetAccountGroupMembers(ctx, accountGroupID)
}
