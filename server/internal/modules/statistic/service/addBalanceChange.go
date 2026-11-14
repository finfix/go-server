package service

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// addBalanceChange учитывает изменение остатка счета наличности в своей валюте: до начала периода —
// во входящий остаток, иначе — в бакет
func (c seriesCalculator) addBalanceChange(opening map[uuid.UUID]decimal.Decimal, changes map[uuid.UUID][]decimal.Decimal, index int, accountID uuid.UUID, amount decimal.Decimal, isMoney bool) {
	if !isMoney {
		return
	}
	if index < 0 {
		opening[accountID] = opening[accountID].Add(amount)
		return
	}
	if changes[accountID] == nil {
		changes[accountID] = make([]decimal.Decimal, len(c.buckets))
	}
	changes[accountID][index] = changes[accountID][index].Add(amount)
}
