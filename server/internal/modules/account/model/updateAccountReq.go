package model

import (
	"pkg/pointer"
	"server/internal/utils/errors"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	repoModel "server/internal/modules/account/repository/model"
	"server/internal/utils/necessary"
)

type UpdateAccountReq struct {
	Necessary          necessary.NecessaryUserInformation
	ID                 uuid.UUID              `json:"id" validate:"required" minimum:"1"` // Идентификатор счета
	Remainder          *decimal.Decimal       `json:"remainder"`                          // Остаток средств на счету
	Name               *string                `json:"name"`                               // Название счета
	IconID             *uuid.UUID             `json:"iconID" minimum:"1"`                 // Идентификатор иконки
	Visible            *bool                  `json:"visible"`                            // Видимость счета
	AccountingInHeader *bool                  `json:"accountingInHeader"`                 // Будет ли счет учитываться в статистике
	AccountingInCharts *bool                  `json:"accountingInCharts"`                 // Будет ли счет учитываться в графиках
	Currency           *string                `json:"currencyCode"`                       // Валюта счета
	SerialNumber       *uint32                `json:"serialNumber"`                       // Порядковый номер счета
	ParentAccountID    *uuid.UUID             `json:"parentAccountID"`                    // Идентификатор родительского счета
	Budget             UpdateAccountBudgetReq `json:"budget"`                             // Месячный бюджет
}

func (s *UpdateAccountReq) ConvertToRepoReq() repoModel.UpdateAccountReq {
	return repoModel.UpdateAccountReq{
		Remainder:          s.Remainder,
		Name:               s.Name,
		IconID:             s.IconID,
		Visible:            s.Visible,
		AccountingInHeader: s.AccountingInHeader,
		AccountingInCharts: s.AccountingInCharts,
		Currency:           s.Currency,
		ParentAccountID:    s.ParentAccountID,
		Budget:             s.Budget.ConvertToRepoReq(),
		SerialNumber:       s.SerialNumber,
	}
}

type UpdateAccountBudgetReq struct {
	Amount         *decimal.Decimal `json:"amount"`         // Сумма
	FixedSum       *decimal.Decimal `json:"fixedSum"`       // Фиксированная сумма
	DaysOffset     *uint32          `json:"daysOffset"`     // Смещение в днях
	GradualFilling *bool            `json:"gradualFilling"` // Постепенное пополнение
}

func (s *UpdateAccountBudgetReq) ConvertToRepoReq() repoModel.UpdateAccountBudgetReq {
	return repoModel.UpdateAccountBudgetReq{
		Amount:         s.Amount,
		FixedSum:       s.FixedSum,
		DaysOffset:     s.DaysOffset,
		GradualFilling: s.GradualFilling,
	}
}

type ProtoUpdateAccountBudgetReq struct {
	*proto.UpdateAccountBudgetRequest
}

func (s ProtoUpdateAccountBudgetReq) ConvertToModel() (UpdateAccountBudgetReq, error) {

	var amount *decimal.Decimal
	if s.Amount != nil {
		amount = pointer.Pointer(decimal.NewFromFloat(*s.Amount))
	}

	var fixedSum *decimal.Decimal
	if s.FixedSum != nil {
		fixedSum = pointer.Pointer(decimal.NewFromFloat(*s.FixedSum))
	}

	var daysOffset *uint32
	if s.DaysOffset != nil {
		daysOffset = s.DaysOffset
	}

	var gradualFilling *bool
	if s.GradualFilling != nil {
		gradualFilling = s.GradualFilling
	}

	return UpdateAccountBudgetReq{
		Amount:         amount,
		FixedSum:       fixedSum,
		DaysOffset:     daysOffset,
		GradualFilling: gradualFilling,
	}, nil
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

	// Convert optional budget
	if p.Budget == nil {
		return res, errors.BadRequest.New("Budget is required")
	}

	budget, err := ProtoUpdateAccountBudgetReq{UpdateAccountBudgetRequest: p.Budget}.ConvertToModel()
	if err != nil {
		return res, err
	}

	var remainder *decimal.Decimal
	if p.Remainder != nil {
		remainder = pointer.Pointer(decimal.NewFromFloat(*p.Remainder))
	}

	var serialNumber *uint32
	if p.SerialNumber != nil {
		serialNumber = p.SerialNumber
	}

	return UpdateAccountReq{

		ID:                 id,
		Name:               p.Name,
		AccountingInCharts: p.AccountingInCharts,
		AccountingInHeader: p.AccountingInHeader,
		Currency:           p.Currency,
		IconID:             iconID,
		ParentAccountID:    parentAccountID,
		Remainder:          remainder,
		SerialNumber:       serialNumber,
		Visible:            p.Visible,
		Budget:             budget,
	}, nil
}
