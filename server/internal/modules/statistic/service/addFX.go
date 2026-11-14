package service

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"server/internal/enum/accountType"
	"server/internal/modules/statistic/model"
)

// Ключи рядов графика fx
const (
	keyExchangeLoss = "exchange_loss" // Потери на обмене
	keyRevaluation  = "revaluation"   // Переоценка
	keyCorrections  = "corrections"   // Корректировки остатка
	keyTransfers    = "transfers"     // Переводы за пределы выбранных счетов
)

// addFX раскладывает изменение стоимости наличности на то, что не объясняется доходами и расходами.
// Каждая операция с выбранными счетами наличности оценивается по курсам на свою дату: потеря на
// обмене — рыночная стоимость отданного минус полученного (у дохода и расхода — разница между суммой
// категории и суммой на счете). Переоценка — остаток: изменение стоимости остатков за бакет минус
// все потоки. Поэтому всегда верно: Δналичности = доходы − расходы − потери + корректировки +
// переводы + переоценка
func (c seriesCalculator) addFX(filters model.Filters) error {
	money := c.balanceAccounts(filters)
	loss := c.get(keyExchangeLoss, "Потери на обмене")
	revaluation := c.get(keyRevaluation, "Переоценка")
	corrections := c.get(keyCorrections, "Корректировки")
	transfers := c.get(keyTransfers, "Переводы за пределы выбранных счетов")

	// Остатки в своей валюте: на начало периода и изменения по бакетам; потоки — в валюте ответа
	opening := make(map[uuid.UUID]decimal.Decimal)
	changes := make(map[uuid.UUID][]decimal.Decimal)
	flows := make([]decimal.Decimal, len(c.buckets))
	for i := range flows {
		flows[i] = decimal.Zero
	}
	for _, transaction := range c.data.transactions {
		fromMoney, toMoney := money[transaction.AccountFromID], money[transaction.AccountToID]
		if !fromMoney && !toMoney || !transaction.DateTransaction.Before(c.endOfPeriod) {
			continue
		}

		// Остатки в своей валюте
		index := c.bucketIndex(transaction.DateTransaction)
		c.addBalanceChange(opening, changes, index, transaction.AccountToID, transaction.AmountTo, toMoney)
		c.addBalanceChange(opening, changes, index, transaction.AccountFromID, transaction.AmountFrom.Neg(), fromMoney)
		if index < 0 {
			continue
		}

		// Обе стороны по курсу на дату операции
		at := endOfDay(transaction.DateTransaction)
		accountFrom, accountTo := c.data.accounts[transaction.AccountFromID], c.data.accounts[transaction.AccountToID]
		valueFrom, err := c.data.convert(transaction.AmountFrom, accountFrom.Currency, at)
		if err != nil {
			return err
		}
		valueTo, err := c.data.convert(transaction.AmountTo, accountTo.Currency, at)
		if err != nil {
			return err
		}
		inflow, outflow := decimal.Zero, decimal.Zero
		if toMoney {
			inflow = valueTo
		}
		if fromMoney {
			outflow = valueFrom
		}
		flow := inflow.Sub(outflow)
		flows[index] = flows[index].Add(flow)

		// Вид потока по второй стороне
		otherType := accountTo.Type
		if toMoney {
			otherType = accountFrom.Type
		}
		switch {
		case otherType == accountType.Balancing:
			corrections.Values[index] = corrections.Values[index].Add(flow)
		case (otherType == accountType.Regular || otherType == accountType.Debt) && !(fromMoney && toMoney):
			transfers.Values[index] = transfers.Values[index].Add(flow)
		default:
			// Обмен между выбранными счетами, доход или расход: стоимость отданного минус полученного
			lost := valueFrom.Sub(valueTo)
			if lost.IsZero() {
				continue
			}
			loss.Values[index] = loss.Values[index].Add(lost)
			pair := c.get("exchange_loss:"+accountFrom.Currency+"→"+accountTo.Currency, "Потери "+accountFrom.Currency+"→"+accountTo.Currency)
			pair.Values[index] = pair.Values[index].Add(lost)
		}
	}

	// Переоценка: изменение стоимости остатков за бакет минус все потоки бакета
	running := make(map[uuid.UUID]decimal.Decimal, len(opening))
	for id, amount := range opening {
		running[id] = amount
	}
	previous, err := c.portfolioValue(running, c.buckets[0].Add(-1))
	if err != nil {
		return err
	}
	for i := range c.buckets {
		for id, bucketChanges := range changes {
			running[id] = running[id].Add(bucketChanges[i])
		}
		current, err := c.portfolioValue(running, c.bucketEnd(i))
		if err != nil {
			return err
		}
		revaluation.Values[i] = current.Sub(previous).Sub(flows[i])
		previous = current
	}
	return nil
}
