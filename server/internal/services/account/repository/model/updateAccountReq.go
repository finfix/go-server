package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UpdateAccountReq struct {
	Remainder          *decimal.Decimal
	Name               *string
	IconID             *int64
	Visible            *bool
	AccountingInHeader *bool
	AccountingInCharts *bool
	Currency           *string
	ParentAccountID    *uuid.UUID
	SerialNumber       *uint32
	Budget             UpdateAccountBudgetReq
}

type UpdateAccountBudgetReq struct {
	Amount         *decimal.Decimal
	FixedSum       *decimal.Decimal
	DaysOffset     *uuid.UUID
	GradualFilling *bool
}
