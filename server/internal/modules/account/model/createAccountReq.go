package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"pkg/datetime"
	"server/internal/utils/necessary"

	"server/internal/modules/account/model/accountType"
	repoModel "server/internal/modules/account/repository/model"
)

type CreateAccountReq struct {
	Necessary          necessary.NecessaryUserInformation
	Name               string                 `json:"name" validate:"required"`                                                          // Название счета
	IconID             uuid.UUID              `json:"iconID" validate:"required" minimum:"1"`                                            // Идентификатор иконки
	Type               accountType.Type       `json:"type" validate:"required" enums:"regular,expense,credit,debt,earnings,investments"` // Тип счета
	Currency           string                 `json:"currency" validate:"required"`                                                      // Валюта счета
	AccountGroupID     uuid.UUID              `json:"accountGroupID" validate:"required" minimum:"1"`                                    // Группа счета
	AccountingInHeader *bool                  `json:"accountingInHeader" validate:"required"`                                            // Подсчет суммы счета в статистике
	AccountingInCharts *bool                  `json:"accountingInCharts" validate:"required"`                                            // Учитывать ли счет в графиках
	DatetimeCreate     datetime.Time          `json:"datetimeCreate" validate:"required"`                                                // Дата создания счета
	Remainder          decimal.Decimal        `json:"remainder"`                                                                         // Остаток средств на счету
	Budget             CreateAccountBudgetReq `json:"budget"`                                                                            // Бюджет
	IsParent           *bool                  `json:"isParent" validate:"required"`                                                      // Является ли счет родительским
	ParentAccountID    *uuid.UUID             `json:"parentAccountID"`                                                                   // Идентификатор родительского счета
	Visible            *bool                  `json:"-"`                                                                                 // Видимость счета
}

func (s CreateAccountReq) Validate(ctx context.Context) error {
	return s.Type.Validate(ctx)
}

func (s CreateAccountReq) ConvertToAccount() Account {
	return Account{
		ID:                 uuid.UUID{},
		Remainder:          s.Remainder,
		Name:               s.Name,
		IconID:             s.IconID,
		Type:               s.Type,
		Currency:           s.Currency,
		Visible:            true,
		AccountGroupID:     s.AccountGroupID,
		AccountingInHeader: *s.AccountingInHeader,
		ParentAccountID:    s.ParentAccountID,
		SerialNumber:       0,
		IsParent:           *s.IsParent,
		CreatedByUserID:    s.Necessary.UserID,
		DatetimeCreate:     datetime.Time{Time: time.Now()},
		AccountingInCharts: *s.AccountingInCharts,
		AccountBudget: AccountBudget{
			Amount:         s.Budget.Amount,
			FixedSum:       s.Budget.FixedSum,
			DaysOffset:     s.Budget.DaysOffset,
			GradualFilling: *s.Budget.GradualFilling,
		},
	}
}

// TODO: Переписать
func (s CreateAccountReq) ConvertToRepoReq() repoModel.CreateAccountReq {
	return repoModel.CreateAccountReq{
		Name:               s.Name,
		IconID:             s.IconID,
		Type:               s.Type,
		Currency:           s.Currency,
		AccountGroupID:     s.AccountGroupID,
		AccountingInHeader: *s.AccountingInHeader,
		AccountingInCharts: *s.AccountingInCharts,
		Budget:             s.Budget.ConvertToRepoReq(),
		IsParent:           *s.IsParent,
		Visible:            true,
		ParentAccountID:    s.ParentAccountID,
		UserID:             s.Necessary.UserID,
		DatetimeCreate:     s.DatetimeCreate.Time,
	}
}

type CreateAccountBudgetReq struct {
	Amount         decimal.Decimal `json:"amount"`                             // Сумма
	FixedSum       decimal.Decimal `json:"fixedSum"`                           // Фиксированная сумма
	DaysOffset     uint32          `json:"daysOffset"`                         // Смещение в днях
	GradualFilling *bool           `json:"gradualFilling" validate:"required"` // Постепенное пополнение
}

// TODO: Переписать
func (s *CreateAccountBudgetReq) ConvertToRepoReq() repoModel.CreateReqBudget {
	return repoModel.CreateReqBudget{
		Amount:         s.Amount,
		FixedSum:       s.FixedSum,
		DaysOffset:     s.DaysOffset,
		GradualFilling: *s.GradualFilling,
	}
}
