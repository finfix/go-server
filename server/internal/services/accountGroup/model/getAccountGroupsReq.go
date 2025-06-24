package model

import (
	"server/internal/utils/necessary"
)

type GetAccountGroupsReq struct {
	Necessary       necessary.NecessaryUserInformation
	AccountGroupIDs []uint32 `json:"accountGroupIDs" schema:"accountGroupIDs" minimum:"1"` // Идентификаторы групп счетов
}
