package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"server/internal/modules/account/model/accountType"
)

type CreateAccountReq struct {
	Budget             CreateReqBudget
	Name               string
	Visible            bool
	IconID             uuid.UUID
	Type               accountType.Type
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
