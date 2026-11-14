package service

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"server/internal/enum/accountType"
	"server/internal/modules/statistic/model"
)

// addFlows раскладывает доходы и расходы. groupBy пустой — без разбивки: ряды income и expense
func (c seriesCalculator) addFlows(filter transactionFilter, chartType model.ChartType, groupBy model.GroupBy) error {
	for _, transaction := range c.data.transactions {
		index := c.bucketIndex(transaction.DateTransaction)
		if index < 0 || !filter.match(c.data, transaction) {
			continue
		}

		// Расход — пополнение расходного счета, доход — списание с доходного
		var accountID uuid.UUID
		var amount decimal.Decimal
		var key string
		accountTo, accountFrom := c.data.accounts[transaction.AccountToID], c.data.accounts[transaction.AccountFromID]
		switch {
		case accountTo.Type == accountType.Expense && chartType != model.ChartIncome:
			accountID, amount, key = accountTo.ID, transaction.AmountTo, keyExpense
		case accountFrom.Type == accountType.Earnings && chartType != model.ChartExpense:
			accountID, amount, key = accountFrom.ID, transaction.AmountFrom, keyIncome
		default:
			continue
		}
		converted, err := c.data.convert(amount, c.data.accounts[accountID].Currency)
		if err != nil {
			return err
		}

		// Кладем сумму в ряды
		for _, series := range c.flowSeries(transaction.ID, accountID, key, groupBy) {
			series.Values[index] = series.Values[index].Add(converted)
		}
	}
	return nil
}
