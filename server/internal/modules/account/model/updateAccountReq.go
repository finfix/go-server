package model

import (
	"server/internal/utils/errors"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"

	repoModel "server/internal/modules/account/repository/model"
	"server/internal/utils/necessary"
)

type UpdateAccountReq struct {
	Necessary          necessary.NecessaryUserInformation
	ID                 uuid.UUID  `json:"id" validate:"required" minimum:"1"` // Идентификатор счета
	Name               *string    `json:"name"`                               // Название счета
	IconID             *uuid.UUID `json:"iconID" minimum:"1"`                 // Идентификатор иконки
	Visible            *bool      `json:"visible"`                            // Видимость счета
	AccountingInHeader *bool      `json:"accountingInHeader"`                 // Будет ли счет учитываться в статистике
	AccountingInCharts *bool      `json:"accountingInCharts"`                 // Будет ли счет учитываться в графиках
	Currency           *string    `json:"currencyCode"`                       // Валюта счета
	Rank               *string    `json:"rank"`                               // Ранг для сортировки счетов (лексикографический, задаётся клиентом)
	ParentAccountID    *uuid.UUID `json:"parentAccountID"`                    // Идентификатор родительского счета
	// LinkedAccountID — установка связывает счёт в мост (симметрично, бэкенд не проверяет доступ
	// ко второму счёту — см. пояснение к proto-контракту). Отсутствие в запросе означает "не
	// менять", поэтому разрыв связи требует отдельного флага UnlinkAccount, а не пустого значения.
	LinkedAccountID *uuid.UUID `json:"linkedAccountID"`
	UnlinkAccount   bool       `json:"unlinkAccount"`
}

func (s *UpdateAccountReq) ConvertToRepoReq() repoModel.UpdateAccountReq {
	return repoModel.UpdateAccountReq{
		Name:               s.Name,
		IconID:             s.IconID,
		Visible:            s.Visible,
		AccountingInHeader: s.AccountingInHeader,
		AccountingInCharts: s.AccountingInCharts,
		Currency:           s.Currency,
		ParentAccountID:    s.ParentAccountID,
		Rank:               s.Rank,
		LinkedAccountID:    s.LinkedAccountID,
		UnlinkAccount:      s.UnlinkAccount,
	}
}

// ProtoUpdateAccountReq wrapper for proto request
type ProtoUpdateAccountReq struct {
	*proto.UpdateAccountRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoUpdateAccountReq) ConvertToModel() (UpdateAccountReq, error) {
	var res UpdateAccountReq

	if p.UpdateAccountRequest == nil {
		return res, errors.BadRequest.New("UpdateAccountRequest is required")
	}

	// Parse ID from bytes
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Parse optional IconID
	var iconID *uuid.UUID
	if p.IconID != nil {
		parsedIconID, err := uuid.FromBytes(p.IconID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		iconID = &parsedIconID
	}

	// Parse optional ParentAccountID
	var parentAccountID *uuid.UUID
	if p.ParentAccountID != nil {
		parsedParentAccountID, err := uuid.FromBytes(p.ParentAccountID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		parentAccountID = &parsedParentAccountID
	}

	// Parse optional LinkedAccountID
	var linkedAccountID *uuid.UUID
	if p.LinkedAccountID != nil {
		parsedLinkedAccountID, err := uuid.FromBytes(p.LinkedAccountID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		linkedAccountID = &parsedLinkedAccountID
	}

	return UpdateAccountReq{

		ID:                 id,
		Name:               p.Name,
		AccountingInCharts: p.AccountingInCharts,
		AccountingInHeader: p.AccountingInHeader,
		Currency:           p.Currency,
		IconID:             iconID,
		ParentAccountID:    parentAccountID,
		Rank:               p.Rank,
		Visible:            p.Visible,
		LinkedAccountID:    linkedAccountID,
		UnlinkAccount:      p.UnlinkAccount,
	}, nil
}
