package model

import (
	"server/internal/enum/accountType"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateAccountReq struct {
	ID                 uuid.UUID
	Budget             CreateReqBudget
	Name               string
	Visible            bool
	IconID             uuid.UUID
	Type               accountType.AccountType
	Currency           string
	AccountGroupID     uuid.UUID
	AccountingInHeader bool
	AccountingInCharts bool
	IsParent           bool
	ParentAccountID    *uuid.UUID
	UserID             uuid.UUID
	DatetimeCreate     time.Time
}

type CreateReqBudget struct {
	Amount         decimal.Decimal
	GradualFilling bool
	FixedSum       decimal.Decimal
	DaysOffset     uint32
}
