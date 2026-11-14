package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"server/internal/enum/accountType"
	"server/internal/modules/statistic/model"
)

// addBudgets добавляет к рядам расходных счетов действующий на конец периода месячный бюджет
func (s *StatisticService) addBudgets(ctx context.Context, req model.GetStatisticReq, data dataset, calculator seriesCalculator, dateTo time.Time) error {
	budgets, err := s.accountBudgetService.GetAccountBudgets(ctx, budgetGetReq(req, data, dateTo))
	if err != nil {
		return err
	}

	// Последняя версия бюджета каждого счета, вступившая в силу к концу периода
	latest := make(map[uuid.UUID]int)
	for i, budget := range budgets {
		if budget.EffectiveFrom.After(dateTo.AddDate(0, 0, 1)) {
			continue
		}
		if j, ok := latest[budget.AccountID]; !ok || budget.EffectiveFrom.After(budgets[j].EffectiveFrom) {
			latest[budget.AccountID] = i
		}
	}

	// Суммируем бюджеты в ряды: дочерние счета — в ряд родителя, если они свернуты
	for accountID, i := range latest {
		account, ok := data.accounts[accountID]
		if !ok || account.Type != accountType.Expense {
			continue
		}
		series, ok := calculator.series[calculator.accountKey(accountID).String()]
		if !ok {
			continue
		}
		amount, err := data.convert(budgets[i].Amount, account.Currency)
		if err != nil {
			return err
		}
		if series.MonthlyBudget == nil {
			series.MonthlyBudget = &decimal.Zero
		}
		total := series.MonthlyBudget.Add(amount)
		series.MonthlyBudget = &total
	}
	return nil
}
