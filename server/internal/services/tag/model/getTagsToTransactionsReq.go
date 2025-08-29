package model

import (
	"github.com/google/uuid"

	"server/internal/utils/necessary"
)

type GetTagsToTransactionsReq struct {
	Necessary       necessary.NecessaryUserInformation
	AccountGroupIDs []uuid.UUID `json:"-" schema:"-" minimum:"1"` // Идентификаторы групп счетов
	TransactionIDs  []uuid.UUID `json:"-" schema:"-" minimum:"1"` // Идентификаторы транзакций
}
