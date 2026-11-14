package service

import (
	"github.com/google/uuid"

	"server/internal/modules/statistic/model"
)

// flowSeries возвращает ряды, в которые попадает сумма транзакции
func (c seriesCalculator) flowSeries(transactionID, accountID uuid.UUID, key string, groupBy model.GroupBy) []*model.Series {
	switch groupBy {
	case model.GroupByAccount:
		return []*model.Series{c.accountSeries(accountID)}
	case model.GroupByTag:
		tags := c.data.tagsByTransaction[transactionID]
		if len(tags) == 0 {
			return []*model.Series{c.get(keyUntagged, "")}
		}
		result := make([]*model.Series, 0, len(tags))
		for _, tagID := range tags {
			series := c.get(tagID.String(), "")
			id := tagID
			series.TagID = &id
			result = append(result, series)
		}
		return result
	}
	names := map[string]string{keyIncome: "Доходы", keyExpense: "Расходы"}
	return []*model.Series{c.get(key, names[key])}
}
