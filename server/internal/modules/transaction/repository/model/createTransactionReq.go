package model

import (
	"server/internal/enum/transactionType"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateTransactionReq struct {
	ID                 uuid.UUID
	Type               transactionType.TransactionType
	AmountFrom         decimal.Decimal
	AmountTo           decimal.Decimal
	Note               string
	AccountFromID      uuid.UUID
	AccountToID        uuid.UUID
	DateTransaction    time.Time
	IsExecuted         bool
	CreatedByUserID    uuid.UUID
	DatetimeCreate     time.Time
	AccountingInCharts bool
	AccountGroupID     uuid.UUID
}
