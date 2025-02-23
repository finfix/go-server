package model

import (
	"server/internal/utils/necessary"
)

type GetTagsToTransactionsReq struct {
	Necessary       necessary.NecessaryUserInformation
	AccountGroupIDs []uint32 `json:"-" schema:"-" minimum:"1"` // Идентификаторы групп счетов
	TransactionIDs  []uint32 `json:"-" schema:"-" minimum:"1"` // Идентификаторы транзакций
}
