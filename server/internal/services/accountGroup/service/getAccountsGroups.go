package service

import (
	"context"

	"server/internal/services/accountGroup/model"
)

// GetAccountGroups Возвращает все группы счетов пользователя
func (s *AccountGroupService) GetAccountGroups(ctx context.Context, req model.GetAccountGroupsReq) ([]model.AccountGroup, error) {
	ctx, span := tracer.Start(ctx, "GetAccountGroups")
	defer span.End()

	return s.accountGroupRepository.GetAccountGroups(ctx, req)
}
