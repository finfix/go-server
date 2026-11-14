package service

import (
	"github.com/shopspring/decimal"

	"server/internal/modules/statistic/model"
)

// summarize считает сводку ряда: сумму, среднее, среднее без нулей, минимум и максимум
func summarize(values []decimal.Decimal) model.SeriesSummary {
	summary := model.SeriesSummary{Sum: decimal.Zero, Average: decimal.Zero, AverageNonZero: decimal.Zero, Min: decimal.Zero, Max: decimal.Zero}
	if len(values) == 0 {
		return summary
	}
	summary.Min, summary.Max = values[0], values[0]
	nonZero := 0
	for _, value := range values {
		summary.Sum = summary.Sum.Add(value)
		summary.Min = decimal.Min(summary.Min, value)
		summary.Max = decimal.Max(summary.Max, value)
		if !value.IsZero() {
			nonZero++
		}
	}
	summary.Average = summary.Sum.Div(decimal.NewFromInt(int64(len(values)))).Round(convertPrecision)
	if nonZero != 0 {
		summary.AverageNonZero = summary.Sum.Div(decimal.NewFromInt(int64(nonZero))).Round(convertPrecision)
	}
	return summary
}
