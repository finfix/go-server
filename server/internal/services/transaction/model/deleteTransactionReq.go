package model

import (
	"server/internal/utils/necessary"
)

type DeleteTransactionReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uint32 `json:"id" validate:"required" minimum:"1"` // Идентификатор транзакции
}
