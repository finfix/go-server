package model

import (
	"server/internal/enum/accountType"
	"time"

	"github.com/google/uuid"
)

type CalculateRemaindersAccountsReq struct {
	IDs             []uuid.UUID
	AccountGroupIDs []uuid.UUID
	Types           []accountType.AccountType
	DateFrom        *time.Time
	DateTo          *time.Time
}
