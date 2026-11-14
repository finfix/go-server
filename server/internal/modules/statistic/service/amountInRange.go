package service

import (
	"github.com/shopspring/decimal"
)

// amountInRange проверяет сумму по границам фильтра
func (f transactionFilter) amountInRange(amount decimal.Decimal) bool {
	if f.filters.AmountMin != nil && amount.LessThan(*f.filters.AmountMin) {
		return false
	}
	if f.filters.AmountMax != nil && amount.GreaterThan(*f.filters.AmountMax) {
		return false
	}
	return true
}
