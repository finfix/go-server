package utils

import (
	"context"

	"server/internal/utils/errors"

	accountModel "server/internal/modules/account/model"
)

func CheckAccountPermissionsForUpdate(ctx context.Context, req accountModel.UpdateAccountReq, permissions accountModel.AccountPermissions) error {

	if req.Currency != nil && !permissions.UpdateCurrency {
		return errors.Forbidden.New("Нельзя менять валюту").WithContextParams(ctx)
	}

	return nil
}
