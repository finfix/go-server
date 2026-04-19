package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/shopspring/decimal"

	"server/internal/modules/settings/repository/currencyDDL"
)

// UpdateCurrencies обновляет курсы валют в базе данных
func (r *SettingsRepository) UpdateCurrencies(ctx context.Context, rates map[string]decimal.Decimal) error {

	q := sq.Insert(currencyDDL.Table).
		Columns(currencyDDL.ColumnSlug, currencyDDL.ColumnName, currencyDDL.ColumnRate, currencyDDL.ColumnSymbol)

	// Формируем аргументы для запроса
	for currency, rate := range rates {
		q = q.Values(currency, currency, rate, currency)
	}

	q = q.Suffix(fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s = EXCLUDED.%s", currencyDDL.ColumnSlug, currencyDDL.ColumnRate, currencyDDL.ColumnRate))

	// Обновляем курсы валют
	return r.db.Exec(ctx, q)
}
