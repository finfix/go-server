package model

import (
	"github.com/google/uuid"
)

type UpdateAccountReq struct {
	Name               *string
	IconID             *uuid.UUID
	Visible            *bool
	AccountingInHeader *bool
	AccountingInCharts *bool
	Currency           *string
	ParentAccountID    *uuid.UUID
	Rank               *string
}
