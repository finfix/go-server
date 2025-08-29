package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	settingsModel "server/internal/modules/settings/model"
	"server/internal/modules/settings/repository/currencyDDL"
)

func (r *SettingsRepository) GetCurrencies(ctx context.Context) (currencies []settingsModel.Currency, err error) {
	ctx, span := tracer.Start(ctx, "GetCurrencies")
	defer span.End()

	return currencies, r.db.Select(ctx, &currencies, sq.
		Select(ddlHelper.SelectAll).
		From(currencyDDL.Table),
	)
}
