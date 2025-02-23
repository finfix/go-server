package model

import (
	"server/internal/utils/necessary"
)

type DeleteAccountGroupReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uint32 `json:"id" schema:"id" validate:"required" minimum:"1"` // Идентификатор счета
}
