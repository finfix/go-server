package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"pkg/datetime"

	"server/internal/modules/transaction/model/transactionType"
)

type CreateTransactionReq struct {
	Type               transactionType.Type
	AmountFrom         decimal.Decimal
	AmountTo           decimal.Decimal
	Note               string
	AccountFromID      uuid.UUID
	AccountToID        uuid.UUID
	DateTransaction    datetime.Date
	IsExecuted         bool
	CreatedByUserID    uuid.UUID
	DatetimeCreate     time.Time
	AccountingInCharts bool
	AccountGroupID     uuid.UUID
}
