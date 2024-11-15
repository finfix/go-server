package service

import (
	"context"

	accountModel "server/internal/services/account/model"
	"server/internal/services/accountPermissions/model"
	"server/internal/services/accountPermissions/service/utils"
)

func (s *AccountPermissionsService) GetAccountsPermissions(ctx context.Context, accounts ...accountModel.Account) (permissions []model.AccountPermissions, err error) {
	ctx, span := tracer.Start(ctx, "GetAccountsPermissions")
	defer span.End()

	// Получаем разрешения из репозитория
	permissionsSet, err := s.accountPermissionsRepository.GetAccountPermissions(ctx)
	if err != nil {
		return permissions, err
	}

	// Собираем разрешения для каждого счета
	permissionsArr := make([]model.AccountPermissions, 0, len(accounts))

	// Проходимся по каждому счету и собираем его разрешения
	for _, account := range accounts {
		permissionsArr = append(permissionsArr, utils.JoinAccountPermissions(
			permissionsSet.TypeToPermissions[account.Type],
			permissionsSet.IsParentToPermissions[account.IsParent],
		))
	}
	return permissionsArr, nil
}
