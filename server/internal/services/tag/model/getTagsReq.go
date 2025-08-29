package model

import (
	"github.com/google/uuid"

	"server/internal/utils/necessary"
)

type GetTagsReq struct {
	Necessary       necessary.NecessaryUserInformation
	AccountGroupIDs []uuid.UUID `json:"accountGroupIDs" schema:"accountGroupIDs" minimum:"1"` // Идентификаторы групп счетов
}
