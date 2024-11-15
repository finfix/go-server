package model

import (
	"pkg/datetime"
	"pkg/necessary"

	repoModel "server/internal/services/accountGroup/repository/model"
)

type CreateAccountGroupReq struct {
	Necessary      necessary.NecessaryUserInformation
	Name           string        `json:"name" db:"name" validate:"required"`                      // Название группы счетов
	Currency       string        `json:"currency" db:"currency_signatura" validate:"required"`    // Валюта группы счетов
	DatetimeCreate datetime.Time `json:"datetimeCreate" db:"datetime_create" validate:"required"` // Дата и время создания группы счетов
}

func (s CreateAccountGroupReq) ConvertToRepoReq() repoModel.CreateAccountGroupReq {
	return repoModel.CreateAccountGroupReq{
		UserID:         s.Necessary.UserID,
		Name:           s.Name,
		Currency:       s.Currency,
		Visible:        true,
		DatetimeCreate: s.DatetimeCreate.Time,
	}
}
