package model

import (
	"context"

	"server/internal/utils/errors"
)

// ChartType — вид графика, как на экране статистики приложения
type ChartType string

const (
	ChartIncomeExpense   ChartType = "income_expense"   // Доходы и расходы двумя рядами
	ChartIncome          ChartType = "income"           // Доходы с разбивкой
	ChartExpense         ChartType = "expense"          // Расходы с разбивкой
	ChartBalanceDetailed ChartType = "balance_detailed" // В наличии по счетам
	ChartBalanceTotal    ChartType = "balance_total"    // В наличии одной линией
	ChartDelta           ChartType = "delta"            // Доходы минус расходы
)

// Validate проверяет, что вид графика известен
func (t ChartType) Validate(ctx context.Context) error {
	switch t {
	case ChartIncomeExpense, ChartIncome, ChartExpense, ChartBalanceDetailed, ChartBalanceTotal, ChartDelta:
		return nil
	}
	return errors.BadRequest.New("unknown chart_type").WithParams("chartType", string(t)).WithContextParams(ctx)
}
