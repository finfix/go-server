package service

import (
	"slices"
	"strings"

	"github.com/google/uuid"

	transactionModel "server/internal/modules/transaction/model"
)

// match проверяет транзакцию по всем фильтрам, кроме дат: их каждый расчет применяет сам
func (f transactionFilter) match(d dataset, transaction transactionModel.Transaction) bool {

	// Типы
	if !f.ignoreTypes && len(f.filters.Types) != 0 && !slices.Contains(f.filters.Types, transaction.Type) {
		return false
	}

	// Заметка
	if f.searchText != "" && !strings.Contains(strings.ToLower(transaction.Note), f.searchText) {
		return false
	}

	// Счета: списание или пополнение
	if len(f.accounts) != 0 && !f.accounts[transaction.AccountFromID] && !f.accounts[transaction.AccountToID] {
		return false
	}
	if f.excludedAccounts[transaction.AccountFromID] || f.excludedAccounts[transaction.AccountToID] {
		return false
	}

	// Валюты счетов
	if len(f.filters.Currencies) != 0 &&
		!slices.Contains(f.filters.Currencies, d.accounts[transaction.AccountFromID].Currency) &&
		!slices.Contains(f.filters.Currencies, d.accounts[transaction.AccountToID].Currency) {
		return false
	}

	// Теги
	tags := d.tagsByTransaction[transaction.ID]
	if len(f.filters.TagIDs) != 0 && !slices.ContainsFunc(tags, func(id uuid.UUID) bool { return slices.Contains(f.filters.TagIDs, id) }) {
		return false
	}
	if slices.ContainsFunc(tags, func(id uuid.UUID) bool { return slices.Contains(f.filters.ExcludedTagIDs, id) }) {
		return false
	}

	// Сумма в валюте счета: подходит, если в диапазоне хотя бы одна сторона
	if f.filters.AmountMin != nil || f.filters.AmountMax != nil {
		if !f.amountInRange(transaction.AmountFrom) && !f.amountInRange(transaction.AmountTo) {
			return false
		}
	}

	return true
}
