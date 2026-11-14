package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/shopspring/decimal"

	"server/internal/modules/settings/repository/currencyRateDDL"
)

// InsertCurrencyRates сохраняет общие курсы провайдера на момент validAt. Повтор того же момента
// перезаписывает курс: загрузку истории можно запускать заново
func (r *SettingsRepository) InsertCurrencyRates(ctx context.Context, validAt time.Time, rates map[string]decimal.Decimal) error {
	if len(rates) == 0 {
		return nil
	}

	q := sq.Insert(currencyRateDDL.Table).
		Columns(currencyRateDDL.ColumnCurrency, currencyRateDDL.ColumnValidAt, currencyRateDDL.ColumnRate, currencyRateDDL.ColumnSource)
	for currency, rate := range rates {
		if !rate.IsPositive() {
			continue
		}
		q = q.Values(currency, validAt, rate, currencyRateDDL.SourceProvider)
	}
	q = q.Suffix("ON CONFLICT (" + currencyRateDDL.ColumnCurrency + ", " + currencyRateDDL.ColumnValidAt + ") WHERE " +
		currencyRateDDL.ColumnAccountGroupID + " IS NULL DO UPDATE SET " + currencyRateDDL.ColumnRate + " = EXCLUDED." + currencyRateDDL.ColumnRate)

	return r.db.Exec(ctx, q)
}
