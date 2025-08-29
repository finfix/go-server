package model

import (
	"github.com/google/uuid"

	"server/internal/utils/necessary"
)

type UpdateTagReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID  `json:"id" validate:"required" minimum:"1"` // Идентификатор подкатегории
	Name      *string `json:"name" minimum:"1"`                   // Название подкатегории
}
