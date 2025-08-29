package utils

import (
	"context"

	"server/internal/utils/errors"

	accountModel "server/internal/modules/account/model"
	accountPermissionsModel "server/internal/modules/accountPermissions/model"
)

func CheckAccountPermissionsForUpdate(ctx context.Context, req accountModel.UpdateAccountReq, permissions accountPermissionsModel.AccountPermissions) error {

	if (req.Budget.DaysOffset != nil || req.Budget.Amount != nil || req.Budget.FixedSum != nil || req.Budget.GradualFilling != nil) && !permissions.UpdateBudget {
		return errors.Forbidden.New("Нельзя менять бюджет").WithContextParams(ctx)
	}

	if req.Currency != nil && !permissions.UpdateCurrency {
		return errors.Forbidden.New("Нельзя менять валюту").WithContextParams(ctx)
	}

	if req.Remainder != nil && !permissions.UpdateRemainder {
		return errors.Forbidden.New("Нельзя менять остаток").WithContextParams(ctx)
	}

	return nil
}
