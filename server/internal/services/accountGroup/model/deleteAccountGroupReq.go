package model

import (
	"github.com/google/uuid"

	"server/internal/utils/necessary"
)

type DeleteAccountGroupReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" schema:"id" validate:"required" minimum:"1"` // Идентификатор счета
}
