package model

import (
	"pkg/datetime"
	"server/internal/enum/accountType"

	"github.com/google/uuid"
)

type CalculateRemaindersAccountsReq struct {
	IDs             []uuid.UUID
	AccountGroupIDs []uuid.UUID
	Types           []accountType.AccountType
	DateFrom        *datetime.Date
	DateTo          *datetime.Date
}
