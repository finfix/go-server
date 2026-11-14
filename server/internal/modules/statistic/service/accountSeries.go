package service

import (
	"github.com/google/uuid"

	"server/internal/modules/statistic/model"
)

// accountSeries возвращает ряд счета
func (c seriesCalculator) accountSeries(accountID uuid.UUID) *model.Series {
	key := c.accountKey(accountID)
	series := c.get(key.String(), c.data.accounts[key].Name)
	series.AccountID = &key
	return series
}
