package service

import (
	"github.com/google/uuid"
)

// withChildren возвращает множество счетов вместе с их дочерними
func (d dataset) withChildren(accountIDs []uuid.UUID) map[uuid.UUID]bool {
	set := make(map[uuid.UUID]bool, len(accountIDs))
	for _, id := range accountIDs {
		set[id] = true
	}
	for id, account := range d.accounts {
		if account.ParentAccountID != nil && set[*account.ParentAccountID] {
			set[id] = true
		}
	}
	return set
}
