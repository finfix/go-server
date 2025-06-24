package model

import (
	"server/internal/utils/necessary"
)

type DeleteTagReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uint32 `json:"id" validate:"required" minimum:"1"` // Идентификатор подкатегории
}
