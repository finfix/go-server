package service

import (
	"strings"

	"github.com/google/uuid"

	"server/internal/modules/statistic/model"
)

// transactionFilter — фильтры с раскрытыми дочерними счетами, готовые к проверке транзакций
type transactionFilter struct {
	filters          model.Filters
	accounts         map[uuid.UUID]bool // Счета фильтра вместе с дочерними
	excludedAccounts map[uuid.UUID]bool // Исключенные счета вместе с дочерними
	searchText       string             // Подстрока заметки в нижнем регистре
	ignoreTypes      bool               // Не применять фильтр типов (график, как в приложении)
}

// newTransactionFilter раскрывает счета фильтров до дочерних
func (d dataset) newTransactionFilter(filters model.Filters, ignoreTypes bool) transactionFilter {
	return transactionFilter{
		filters:          filters,
		accounts:         d.withChildren(filters.AccountIDs),
		excludedAccounts: d.withChildren(filters.ExcludedAccountIDs),
		searchText:       strings.ToLower(filters.SearchText),
		ignoreTypes:      ignoreTypes,
	}
}
