package model

import (
	"pkg/datetime"

	"server/internal/services/account/model/accountType"

	"github.com/google/uuid"
)

type CalculateRemaindersAccountsReq struct {
	IDs             []uuid.UUID
	AccountGroupIDs []uuid.UUID
	Types           []accountType.Type
	DateFrom        *datetime.Date
	DateTo          *datetime.Date
}
