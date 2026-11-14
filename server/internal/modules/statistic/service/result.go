package service

import (
	"github.com/shopspring/decimal"

	"server/internal/modules/statistic/model"
)

// result собирает ряды со сводками в стабильном порядке
func (c seriesCalculator) result(chartType model.ChartType) []model.Series {

	// Дельта — один ряд из доходов и расходов
	if chartType == model.ChartDelta {
		income, expense := c.get(keyIncome, "Доходы"), c.get(keyExpense, "Расходы")
		delta := model.Series{Key: keyDelta, Name: "Доходы минус расходы", Values: make([]decimal.Decimal, len(c.buckets))} //nolint:exhaustruct
		for i := range delta.Values {
			delta.Values[i] = income.Values[i].Sub(expense.Values[i])
		}
		delta.Summary = summarize(delta.Values)
		return []model.Series{delta}
	}

	// Доходы и расходы — оба ряда всегда, даже пустые
	if chartType == model.ChartIncomeExpense {
		c.get(keyIncome, "Доходы")
		c.get(keyExpense, "Расходы")
	}

	result := make([]model.Series, 0, len(c.series))
	for _, series := range c.series {
		series.Summary = summarize(series.Values)
		result = append(result, *series)
	}
	sortSeries(result)
	return result
}
