package model

import (
	"pkg/datetime"
	"server/internal/utils/necessary"

	repoModel "server/internal/services/tag/repository/model"
)

type CreateTagReq struct {
	Necessary      necessary.NecessaryUserInformation
	AccountGroupID uint32        `json:"accountGroupID" validate:"required"` // Идентификатор группы счетов
	Name           string        `json:"name" validate:"required"`           // Название подкатегории
	DatetimeCreate datetime.Time `json:"datetimeCreate" validate:"required"` // Дата создания подкатегории
}

func (s CreateTagReq) ConvertToRepoReq() repoModel.CreateTagReq {
	return repoModel.CreateTagReq{
		Name:            s.Name,
		AccountGroupID:  s.AccountGroupID,
		CreatedByUserID: s.Necessary.UserID,
		DatetimeCreate:  s.DatetimeCreate.Time,
	}
}
