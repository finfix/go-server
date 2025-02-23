package utils

import (
	"context"

	"pkg/errors"

	accountModel "server/internal/services/account/model"
	accountPermissionsModel "server/internal/services/accountPermissions/model"
)

func CheckAccountPermissionsForUpdate(ctx context.Context, req accountModel.UpdateAccountReq, permissions accountPermissionsModel.AccountPermissions) error {

	if (req.Budget.DaysOffset != nil || req.Budget.Amount != nil || req.Budget.FixedSum != nil || req.Budget.GradualFilling != nil) && !permissions.UpdateBudget {
		return errors.Forbidden.New(ctx, "Нельзя менять бюджет")
	}

	if req.Currency != nil && !permissions.UpdateCurrency {
		return errors.Forbidden.New(ctx, "Нельзя менять валюту")
	}

	if req.Remainder != nil && !permissions.UpdateRemainder {
		return errors.Forbidden.New(ctx, "Нельзя менять остаток")
	}

	return nil
}
