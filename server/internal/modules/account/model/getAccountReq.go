package model

import (
	"context"
	"server/internal/enum/accountType"
	"server/internal/utils/errors"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"

	"pkg/datetime"
	"server/internal/utils/necessary"
	repoModel "server/internal/modules/account/repository/model"
)

type GetAccountsReq struct {
	Necessary          necessary.NecessaryUserInformation
	Type               *accountType.AccountType `json:"type" schema:"type" enums:"regular,expense,credit,debt,earnings,investments"` // Тип счета
	AccountingInHeader *bool                    `json:"accountingInHeader" schema:"accountingInHeader"`                              // Учитывать ли счет в шапке
	AccountingInCharts *bool                    `json:"accountingInCharts" schema:"accountingInCharts"`                              // Учитывать ли счет в графиках
	AccountGroupIDs    []uuid.UUID              `json:"accountGroupIDs" schema:"accountGroupIDs" minimum:"1"`                        // Идентификаторы групп счетов
	DateFrom           *datetime.Date           `json:"dateFrom" schema:"dateFrom" format:"date" swaggertype:"primitive,string"`     // Дата начала выборки (Обязательна при type = expense or earnings и отсутствующем периоде)
	DateTo             *datetime.Date           `json:"dateTo" schema:"dateTo" format:"date" swaggertype:"primitive,string"`         // Дата конца выборки (Обязательна при type = expense or earnings и отсутствующем периоде)
	Visible            *bool                    `json:"visible" schema:"visible"`                                                    // Видимость счета
	Currency           *string                  `json:"-" schema:"-"`                                                                // Валюта счета
	IsParent           *bool                    `json:"-" schema:"-"`                                                                // Является ли счет родительским
	IDs                []uuid.UUID              `json:"-" schema:"ids"`
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
		req.Types = []accountType.AccountType{*s.Type}
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

// ProtoGetAccountsReq wrapper for proto request
type ProtoGetAccountsReq struct {
	*proto.GetAccountsRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoGetAccountsReq) ConvertToModel() (GetAccountsReq, error) {
	var res GetAccountsReq

	if p.GetAccountsRequest == nil {
		return res, errors.BadRequest.New("GetAccountsRequest is required")
	}

	// Convert account group IDs
	accountGroupIDs := make([]uuid.UUID, 0, len(p.AccountGroupIDs))
	for _, idBytes := range p.AccountGroupIDs {
		id, err := uuid.FromBytes(idBytes)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		accountGroupIDs = append(accountGroupIDs, id)
	}

	// Convert account type if provided
	var _accountType *accountType.AccountType
	if p.Type != nil {
		__accountType, err := accountType.ProtoAccountType{AccountType: *p.Type}.ConvertToModel()
		if err != nil {
			return res, err
		}
		_accountType = &__accountType
	}

	var dateFrom *datetime.Date
	if p.DateFrom != nil {
		dateFrom = &datetime.Date{Time: p.DateFrom.AsTime()}
	}

	var dateTo *datetime.Date
	if p.DateTo != nil {
		dateTo = &datetime.Date{Time: p.DateTo.AsTime()}
	}

	return GetAccountsReq{
		AccountGroupIDs:    accountGroupIDs,
		AccountingInCharts: p.AccountingInCharts,
		AccountingInHeader: p.AccountingInHeader,
		DateFrom:           dateFrom,
		DateTo:             dateTo,
		Type:               _accountType,
		Visible:            p.Visible,
	}, nil
}
