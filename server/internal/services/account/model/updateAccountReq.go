package model

import (
	"github.com/shopspring/decimal"

	repoModel "server/internal/services/account/repository/model"
	"server/internal/utils/necessary"
)

type UpdateAccountReq struct {
	Necessary          necessary.NecessaryUserInformation
	ID                 uint32                 `json:"id" validate:"required" minimum:"1"` // Идентификатор счета
	Remainder          *decimal.Decimal       `json:"remainder"`                          // Остаток средств на счету
	Name               *string                `json:"name"`                               // Название счета
	IconID             *uint32                `json:"iconID" minimum:"1"`                 // Идентификатор иконки
	Visible            *bool                  `json:"visible"`                            // Видимость счета
	AccountingInHeader *bool                  `json:"accountingInHeader"`                 // Будет ли счет учитываться в статистике
	AccountingInCharts *bool                  `json:"accountingInCharts"`                 // Будет ли счет учитываться в графиках
	Currency           *string                `json:"currencyCode"`                       // Валюта счета
	SerialNumber       *uint32                `json:"serialNumber"`                       // Порядковый номер счета
	ParentAccountID    *uint32                `json:"parentAccountID"`                    // Идентификатор родительского счета
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
