package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"server/internal/modules/statistic/model"
)

// addBalances собирает остатки счетов regular и debt на конец каждого бакета от первой транзакции.
// Из фильтров учитываются счета, исключенные счета и валюты — они выбирают, какие счета показать
func (c seriesCalculator) addBalances(filters model.Filters, total bool) error {
	// Счета наличности по фильтрам
	balanceAccounts := c.balanceAccounts(filters)

	// Изменения остатков по бакетам в валюте счета; всё до начала периода — во входящий остаток
	changes := make(map[uuid.UUID][]decimal.Decimal)
	opening := make(map[uuid.UUID]decimal.Decimal)
	apply := func(accountID uuid.UUID, amount decimal.Decimal, date time.Time) {
		if !balanceAccounts[accountID] || !date.Before(c.endOfPeriod) {
			return
		}
		if _, ok := changes[accountID]; !ok {
			changes[accountID] = make([]decimal.Decimal, len(c.buckets))
		}
		if index := c.bucketIndex(date); index >= 0 {
			changes[accountID][index] = changes[accountID][index].Add(amount)
		} else {
			opening[accountID] = opening[accountID].Add(amount)
		}
	}
	for _, transaction := range c.data.transactions {
		apply(transaction.AccountToID, transaction.AmountTo, transaction.DateTransaction)
		apply(transaction.AccountFromID, transaction.AmountFrom.Neg(), transaction.DateTransaction)
	}

	// Накопительный остаток каждого счета переводим в валюту ответа по курсу на конец бакета
	for accountID := range balanceAccounts {
		var series *model.Series
		if total {
			series = c.get(keyBalance, "В наличии")
		} else {
			series = c.accountSeries(accountID)
		}
		running := opening[accountID]
		for i := range c.buckets {
			if changes[accountID] != nil {
				running = running.Add(changes[accountID][i])
			}
			value, err := c.data.convert(running, c.data.accounts[accountID].Currency, c.bucketEnd(i))
			if err != nil {
				return err
			}
			series.Values[i] = series.Values[i].Add(value)
		}
	}
	return nil
}
