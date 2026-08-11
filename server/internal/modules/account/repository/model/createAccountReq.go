package model

import (
	"server/internal/enum/accountType"
	"time"

	"github.com/google/uuid"
)

type CreateAccountReq struct {
	ID                 uuid.UUID
	Name               string
	Visible            bool
	IconID             uuid.UUID
	Type               accountType.AccountType
	Currency           string
	AccountGroupID     uuid.UUID
	AccountingInHeader bool
	AccountingInCharts bool
	Rank               string
	IsParent           bool
	ParentAccountID    *uuid.UUID
	UserID             uuid.UUID
	DatetimeCreate     time.Time
}
