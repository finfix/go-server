package service

import (
	"slices"

	"github.com/google/uuid"

	"server/internal/enum/accountType"
	"server/internal/modules/statistic/model"
)

// balanceAccounts — счета наличности (regular и debt) по фильтрам счетов, исключенных счетов и валют
func (c seriesCalculator) balanceAccounts(filters model.Filters) map[uuid.UUID]bool {
	filter := c.data.newTransactionFilter(filters, true)
	result := make(map[uuid.UUID]bool)
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
		result[id] = true
	}
	return result
}
