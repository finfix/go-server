package utils

import (
	"server/internal/modules/accountPermissions/model"
)

func JoinAccountPermissions(permissions ...model.AccountPermissions) (joinedPermissions model.AccountPermissions) {
	joinedPermissions = allTruePermissions
	for _, permission := range permissions {
		joinedPermissions.UpdateBudget = joinedPermissions.UpdateBudget && permission.UpdateBudget
		joinedPermissions.UpdateRemainder = joinedPermissions.UpdateRemainder && permission.UpdateRemainder
		joinedPermissions.UpdateCurrency = joinedPermissions.UpdateCurrency && permission.UpdateCurrency
		joinedPermissions.UpdateParentAccountID = joinedPermissions.UpdateParentAccountID && permission.UpdateParentAccountID

		joinedPermissions.CreateTransaction = joinedPermissions.CreateTransaction && permission.CreateTransaction
	}
	return joinedPermissions
}

var allTruePermissions = model.AccountPermissions{
	UpdateBudget:          true,
	UpdateRemainder:       true,
	UpdateCurrency:        true,
	UpdateParentAccountID: true,

	CreateTransaction: true,
}
