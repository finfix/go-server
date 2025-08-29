package model

import (
	"github.com/google/uuid"

	"server/internal/modules/account/model/accountType"
)

type GetAccountsReq struct {
	IDs                []uuid.UUID
	AccountGroupIDs    []uuid.UUID
	Types              []accountType.Type
	AccountingInHeader *bool
	AccountingInCharts *bool
	Visible            *bool
	Currencies         []string
	IsParent           *bool
	ParentAccountIDs   []uuid.UUID
}
