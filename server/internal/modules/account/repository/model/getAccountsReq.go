package model

import (
	"server/internal/enum/accountType"

	"github.com/google/uuid"
)

type GetAccountsReq struct {
	IDs                []uuid.UUID
	AccountGroupIDs    []uuid.UUID
	Types              []accountType.AccountType
	AccountingInHeader *bool
	AccountingInCharts *bool
	Visible            *bool
	Currencies         []string
	IsParent           *bool
	ParentAccountIDs   []uuid.UUID
}
