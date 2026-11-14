package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	settingsModel "server/internal/modules/settings/model"
	"server/internal/modules/settings/repository/currencyRateDDL"
)

// GetCurrencyRates возвращает общие точки истории курсов по валютам и периоду, по возрастанию момента (при равенстве — по коду валюты)
func (r *SettingsRepository) GetCurrencyRates(ctx context.Context, req settingsModel.GetCurrencyRatesReq) ([]settingsModel.CurrencyRate, error) {
	ctx, span := tracer.Start(ctx, "GetCurrencyRates")
	defer span.End()

	q := sq.Select(currencyRateDDL.ColumnCurrency, currencyRateDDL.ColumnValidAt, currencyRateDDL.ColumnRate).
		From(currencyRateDDL.Table).
		Where(sq.Eq{currencyRateDDL.ColumnAccountGroupID: nil}).
		OrderBy(currencyRateDDL.ColumnValidAt, currencyRateDDL.ColumnCurrency)
	if len(req.Currencies) != 0 {
		q = q.Where(sq.Eq{currencyRateDDL.ColumnCurrency: req.Currencies})
	}
	if req.DateFrom != nil {
		q = q.Where(sq.GtOrEq{currencyRateDDL.ColumnValidAt: *req.DateFrom})
	}
	if req.DateTo != nil {
		q = q.Where(sq.Lt{currencyRateDDL.ColumnValidAt: req.DateTo.AddDate(0, 0, 1)})
	}

	if req.Since != nil {
		q = q.Where(sq.Gt{currencyRateDDL.ColumnValidAt: *req.Since})
	}
	if req.Limit > 0 {
		q = q.Limit(uint64(req.Limit))
	}

	var rates []settingsModel.CurrencyRate
	if err := r.db.Select(ctx, &rates, q); err != nil {
		return nil, err
	}
	return rates, nil
}
