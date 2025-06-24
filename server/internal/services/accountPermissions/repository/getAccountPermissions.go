package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"pkg/log"
	permissionAccountType "server/internal/services/accountPermissions/repository/accountType"
	"server/internal/services/accountPermissions/repository/actionType"

	"server/internal/services/account/model/accountType"
	"server/internal/services/accountPermissions/model"
)

type permissionItem struct {
	AccountType permissionAccountType.Type `db:"account_type"`
	ActionType  actionType.Type            `db:"action_type"`
	Access      bool                       `db:"access"`
}

func (r *AccountPermissionsRepository) GetAccountPermissions(ctx context.Context) (permissionSet model.PermissionSet, err error) {
	ctx, span := tracer.Start(ctx, "GetAccountPermissions")
	defer span.End()

	// Проверяем наличие данных в кэше
	var ok bool
	if permissionSet, ok = r.cache.Get(struct{}{}); ok {
		return permissionSet, nil
	}

	log.WithContextParams(ctx).Info("Обновляем кэш с пермишенами на действия со счетами")

	var permissions []permissionItem

	// Делаем запрос в базу
	if err = r.db.Select(ctx, &permissions, sq.
		Select(ddlHelper.SelectAll).
		From(TableName),
	); err != nil {
		return permissionSet, err
	}

	// Инициализируем структуру
	permissionSet = model.PermissionSet{
		TypeToPermissions:     make(map[accountType.Type]model.AccountPermissions),
		IsParentToPermissions: make(map[bool]model.AccountPermissions),
	}

	// Проходимся по каждой строке из базы данных
	for _, permissionItem := range permissions {

		accountTypeClassification := permissionAccountType.ClassificationMatching[permissionItem.AccountType]

		// Получаем права доступа, которые уже лежат в объекте permissionSet для данного типа аккаунта
		var permission model.AccountPermissions
		switch accountTypeClassification {
		case permissionAccountType.AccountTypeByType:
			permission = permissionSet.TypeToPermissions[accountType.Type(permissionItem.AccountType)]
		case permissionAccountType.AccountTypeByParent:
			permission = permissionSet.IsParentToPermissions[permissionItem.AccountType == permissionAccountType.Parent]
		}

		// Смотрим на действие и присваиваем соответствующий доступ
		switch permissionItem.ActionType {
		case actionType.UpdateBudget:
			permission.UpdateBudget = permissionItem.Access
		case actionType.UpdateRemainder:
			permission.UpdateRemainder = permissionItem.Access
		case actionType.UpdateCurrency:
			permission.UpdateCurrency = permissionItem.Access
		case actionType.UpdateParentAccountID:
			permission.UpdateParentAccountID = permissionItem.Access
		case actionType.CreateTransaction:
			permission.CreateTransaction = permissionItem.Access
		}

		// Сохраняем права доступа в объект permissionSet до следующей итерации
		switch accountTypeClassification {
		case permissionAccountType.AccountTypeByType:
			permissionSet.TypeToPermissions[accountType.Type(permissionItem.AccountType)] = permission
		case permissionAccountType.AccountTypeByParent:
			permissionSet.IsParentToPermissions[permissionItem.AccountType == permissionAccountType.Parent] = permission
		}
	}

	// Сохраняем данные в кэш
	r.cache.Set(struct{}{}, permissionSet)

	// Возвращаем данные
	return permissionSet, nil
}
