package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/settings/repository/currencyRateDDL"
)

// GetProviderRateDays возвращает дни (UTC), за которые уже есть общие курсы провайдера
func (r *SettingsRepository) GetProviderRateDays(ctx context.Context, from, to time.Time) (map[string]bool, error) {
	var days []string
	err := r.db.Select(ctx, &days, sq.Select("DISTINCT to_char("+currencyRateDDL.ColumnValidAt+" AT TIME ZONE 'UTC', 'YYYY-MM-DD')").
		From(currencyRateDDL.Table).
		Where(sq.Eq{currencyRateDDL.ColumnSource: currencyRateDDL.SourceProvider, currencyRateDDL.ColumnAccountGroupID: nil}).
		Where(sq.GtOrEq{currencyRateDDL.ColumnValidAt: from}).
		Where(sq.Lt{currencyRateDDL.ColumnValidAt: to}))
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool, len(days))
	for _, day := range days {
		result[day] = true
	}
	return result, nil
}
