package service

import (
	transactionModel "server/internal/modules/transaction/model"
)

// inDates проверяет дату транзакции по границам фильтра включительно
func (f transactionFilter) inDates(transaction transactionModel.Transaction) bool {
	if f.filters.DateFrom != nil && transaction.DateTransaction.Before(*f.filters.DateFrom) {
		return false
	}
	if f.filters.DateTo != nil && transaction.DateTransaction.After(*f.filters.DateTo) {
		return false
	}
	return true
}
