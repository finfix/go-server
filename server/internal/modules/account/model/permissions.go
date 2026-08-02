package model

import (
	"server/internal/enum/accountType"
	"server/internal/utils/errors"
)

// AccountPermissions - разрешения на действия со счетом
type AccountPermissions struct {
	UpdateBudget          bool
	UpdateRemainder       bool
	UpdateCurrency        bool
	UpdateParentAccountID bool

	CreateTransaction bool
}

// typeToPermissions - разрешения на действия со счетом в зависимости от его типа
var typeToPermissions = map[accountType.AccountType]AccountPermissions{
	accountType.Regular: {
		UpdateRemainder:       true,
		UpdateCurrency:        true,
		UpdateParentAccountID: true,
		CreateTransaction:     true,
	},
	accountType.Expense: {
		UpdateCurrency:        true,
		UpdateParentAccountID: true,
		UpdateBudget:          true,
		CreateTransaction:     true,
	},
	accountType.Earnings: {
		UpdateCurrency:        true,
		UpdateParentAccountID: true,
		UpdateBudget:          true,
		CreateTransaction:     true,
	},
	accountType.Debt: {
		UpdateBudget:          true,
		UpdateCurrency:        true,
		UpdateParentAccountID: true,
		CreateTransaction:     true,
	},
}

// isParentToPermissions - разрешения на действия со счетом в зависимости от того,
// является ли он родительским. Комбинируются (по И) с разрешениями из typeToPermissions.
var isParentToPermissions = map[bool]AccountPermissions{
	// Обычный (не родительский) счет
	false: {
		UpdateCurrency:        true,
		UpdateParentAccountID: true,
		UpdateBudget:          true,
		UpdateRemainder:       true,
		CreateTransaction:     true,
	},
	// Родительский счет: менять родителя и создавать транзакции нельзя
	true: {
		UpdateRemainder: true,
		UpdateCurrency:  true,
		UpdateBudget:    true,
	},
}

// GetAccountPermissions возвращает разрешения на действия для одного счета
func GetAccountPermissions(account Account) (AccountPermissions, error) {
	permissionsByType, ok := typeToPermissions[account.Type]
	if !ok {
		return AccountPermissions{}, errors.InternalServer.New("Не найдены разрешения для типа счета").
			WithParams("type", account.Type)
	}

	permissionsByIsParent, ok := isParentToPermissions[account.IsParent]
	if !ok {
		return AccountPermissions{}, errors.InternalServer.New("Не найдены разрешения для признака родительского счета").
			WithParams("isParent", account.IsParent)
	}

	return joinAccountPermissions(permissionsByType, permissionsByIsParent), nil
}

// GetAccountsPermissions возвращает разрешения на действия для каждого из переданных счетов
func GetAccountsPermissions(accounts ...Account) ([]AccountPermissions, error) {
	permissions := make([]AccountPermissions, 0, len(accounts))
	for _, account := range accounts {
		permission, err := GetAccountPermissions(account)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func joinAccountPermissions(permissions ...AccountPermissions) AccountPermissions {
	joined := AccountPermissions{
		UpdateBudget:          true,
		UpdateRemainder:       true,
		UpdateCurrency:        true,
		UpdateParentAccountID: true,
		CreateTransaction:     true,
	}
	for _, permission := range permissions {
		joined.UpdateBudget = joined.UpdateBudget && permission.UpdateBudget
		joined.UpdateRemainder = joined.UpdateRemainder && permission.UpdateRemainder
		joined.UpdateCurrency = joined.UpdateCurrency && permission.UpdateCurrency
		joined.UpdateParentAccountID = joined.UpdateParentAccountID && permission.UpdateParentAccountID
		joined.CreateTransaction = joined.CreateTransaction && permission.CreateTransaction
	}
	return joined
}