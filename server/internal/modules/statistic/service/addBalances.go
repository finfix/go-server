package service

import (
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"server/internal/enum/accountType"
	"server/internal/modules/statistic/model"
)

// addBalances собирает остатки счетов regular и debt на конец каждого бакета от первой транзакции.
// Из фильтров учитываются счета, исключенные счета и валюты — они выбирают, какие счета показать
func (c seriesCalculator) addBalances(filters model.Filters, total bool) error {
	filter := c.data.newTransactionFilter(filters, true)

	// Счета наличности по фильтрам
	balanceAccounts := make(map[uuid.UUID]bool)
	for id, account := range c.data.accounts {
		if account.Type != accountType.Regular && account.Type != accountType.Debt {
			continue
		}
		if len(filter.accounts) != 0 && !filter.accounts[id] || filter.excludedAccounts[id] {
			continue
		}
		if len(filters.Currencies) != 0 && !slices.Contains(filters.Currencies, account.Currency) {
			continue
		}
		balanceAccounts[id] = true
	}

	// Изменения остатков по бакетам; всё до начала периода — во входящий остаток
	changes := make(map[uuid.UUID][]decimal.Decimal)
	opening := make(map[uuid.UUID]decimal.Decimal)
	apply := func(accountID uuid.UUID, amount decimal.Decimal, date time.Time) error {
		if !balanceAccounts[accountID] || !date.Before(c.endOfPeriod) {
			return nil
		}
		converted, err := c.data.convert(amount, c.data.accounts[accountID].Currency)
		if err != nil {
			return err
		}
		if _, ok := changes[accountID]; !ok {
			changes[accountID] = make([]decimal.Decimal, len(c.buckets))
		}
		if index := c.bucketIndex(date); index >= 0 {
			changes[accountID][index] = changes[accountID][index].Add(converted)
		} else {
			opening[accountID] = opening[accountID].Add(converted)
		}
		return nil
	}
	for _, transaction := range c.data.transactions {
		if err := apply(transaction.AccountToID, transaction.AmountTo, transaction.DateTransaction); err != nil {
			return err
		}
		if err := apply(transaction.AccountFromID, transaction.AmountFrom.Neg(), transaction.DateTransaction); err != nil {
			return err
		}
	}

	// Накопительный итог по каждому счету и сложение в ряды
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
			series.Values[i] = series.Values[i].Add(running)
		}
	}
	return nil
}
