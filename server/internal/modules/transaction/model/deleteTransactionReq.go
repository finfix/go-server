package model

import (
	"github.com/google/uuid"

	"server/internal/utils/necessary"
)

type DeleteTransactionReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required" minimum:"1"` // Идентификатор транзакции
}
