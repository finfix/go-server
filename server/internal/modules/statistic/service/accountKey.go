package service

import (
	"github.com/google/uuid"
)

// accountKey возвращает счет ряда: родителя, если счет свернут в него
func (c seriesCalculator) accountKey(accountID uuid.UUID) uuid.UUID {
	account := c.data.accounts[accountID]
	if !c.aggregateIntoParents || account.ParentAccountID == nil {
		return accountID
	}
	if parent, ok := c.data.accounts[*account.ParentAccountID]; ok && parent.IsParent {
		return parent.ID
	}
	return accountID
}
