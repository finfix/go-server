package model

import (
	"context"

	"pkg/datetime"
	"server/internal/utils/necessary"

	"server/internal/services/account/model/accountType"
	repoModel "server/internal/services/account/repository/model"
)

type GetAccountsReq struct {
	Necessary          necessary.NecessaryUserInformation
	Type               *accountType.Type `json:"type" schema:"type" enums:"regular,expense,credit,debt,earnings,investments"` // Тип счета
	AccountingInHeader *bool             `json:"accountingInHeader" schema:"accountingInHeader"`                              // Учитывать ли счет в шапке
	AccountingInCharts *bool             `json:"accountingInCharts" schema:"accountingInCharts"`                              // Учитывать ли счет в графиках
	AccountGroupIDs    []uint32          `json:"accountGroupIDs" schema:"accountGroupIDs" minimum:"1"`                        // Идентификаторы групп счетов
	DateFrom           *datetime.Date    `json:"dateFrom" schema:"dateFrom" format:"date" swaggertype:"primitive,string"`     // Дата начала выборки (Обязательна при type = expense or earnings и отсутствующем периоде)
	DateTo             *datetime.Date    `json:"dateTo" schema:"dateTo" format:"date" swaggertype:"primitive,string"`         // Дата конца выборки (Обязательна при type = expense or earnings и отсутствующем периоде)
	Visible            *bool             `json:"visible" schema:"visible"`                                                    // Видимость счета
	Currency           *string           `json:"-" schema:"-"`                                                                // Валюта счета
	IsParent           *bool             `json:"-" schema:"-"`                                                                // Является ли счет родительским
	IDs                []uint32          `json:"-" schema:"-"`
}

func (s GetAccountsReq) Validate(ctx context.Context) error {
	return s.Type.Validate(ctx)
}

// TODO: Переписать
func (s *GetAccountsReq) ConvertToRepoReq() repoModel.GetAccountsReq {
	var req repoModel.GetAccountsReq
	req.IDs = s.IDs
	req.AccountGroupIDs = s.AccountGroupIDs
	if s.Type != nil {
		req.Types = []accountType.Type{*s.Type}
	}
	req.AccountingInHeader = s.AccountingInHeader
	req.AccountingInCharts = s.AccountingInCharts
	req.Visible = s.Visible
	if s.Currency != nil {
		req.Currencies = []string{*s.Currency}
	}
	req.IsParent = s.IsParent

	return req
}
