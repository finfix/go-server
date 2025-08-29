package model

import (
	"context"
	"server/internal/enum/accountType"
	"server/internal/utils/errors"
	"time"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"pkg/datetime"
	"server/internal/utils/necessary"

	repoModel "server/internal/modules/account/repository/model"
)

type CreateAccountReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"` // Идентификатор счета

	Name               string                  `json:"name" validate:"required"`                                                          // Название счета
	IconID             uuid.UUID               `json:"iconID" validate:"required" minimum:"1"`                                            // Идентификатор иконки
	Type               accountType.AccountType `json:"type" validate:"required" enums:"regular,expense,credit,debt,earnings,investments"` // Тип счета
	Currency           string                  `json:"currency" validate:"required"`                                                      // Валюта счета
	AccountGroupID     uuid.UUID               `json:"accountGroupID" validate:"required" minimum:"1"`                                    // Группа счета
	AccountingInHeader bool                    `json:"accountingInHeader"`                                                                // Подсчет суммы счета в статистике
	AccountingInCharts bool                    `json:"accountingInCharts"`                                                                // Учитывать ли счет в графиках
	DatetimeCreate     time.Time               `json:"datetimeCreate" validate:"required"`                                                // Дата создания счета
	Remainder          decimal.Decimal         `json:"remainder"`                                                                         // Остаток средств на счету
	Budget             *CreateAccountBudgetReq `json:"budget"`                                                                            // Бюджет
	IsParent           bool                    `json:"isParent"`                                                                          // Является ли счет родительским
	ParentAccountID    *uuid.UUID              `json:"parentAccountID"`                                                                   // Идентификатор родительского счета
	Visible            bool                    `json:"-"`                                                                                 // Видимость счета
}

func (s CreateAccountReq) Validate(ctx context.Context) error {
	return s.Type.Validate(ctx)
}

func (s CreateAccountReq) ConvertToAccount() Account {
	return Account{
		ID:                 s.ID,
		Remainder:          s.Remainder,
		Name:               s.Name,
		IconID:             s.IconID,
		Type:               s.Type,
		Currency:           s.Currency,
		Visible:            true,
		AccountGroupID:     s.AccountGroupID,
		AccountingInHeader: s.AccountingInHeader,
		ParentAccountID:    s.ParentAccountID,
		SerialNumber:       0,
		IsParent:           s.IsParent,
		CreatedByUserID:    s.Necessary.UserID,
		DatetimeCreate:     datetime.Time{Time: time.Now()},
		AccountingInCharts: s.AccountingInCharts,
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
		ID:                 s.ID,
		Name:               s.Name,
		IconID:             s.IconID,
		Type:               s.Type,
		Currency:           s.Currency,
		AccountGroupID:     s.AccountGroupID,
		AccountingInHeader: s.AccountingInHeader,
		AccountingInCharts: s.AccountingInCharts,
		Budget:             s.Budget.ConvertToRepoReq(),
		IsParent:           s.IsParent,
		Visible:            true,
		ParentAccountID:    s.ParentAccountID,
		UserID:             s.Necessary.UserID,
		DatetimeCreate:     s.DatetimeCreate,
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

// ProtoCreateAccountReq wrapper for proto request
type ProtoCreateAccountReq struct {
	*proto.CreateAccountRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoCreateAccountReq) ConvertToModel() (CreateAccountReq, error) {
	var res CreateAccountReq

	if p.CreateAccountRequest == nil {
		return res, errors.BadRequest.New("CreateAccountRequest is required")
	}

	// Parse ID from bytes
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Parse IconID
	iconID, err := uuid.FromBytes(p.IconID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Parse AccountGroupID
	accountGroupID, err := uuid.FromBytes(p.AccountGroupID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Convert account type
	accountType, err := accountType.ProtoAccountType{AccountType: p.Type}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Convert datetime
	if p.DatetimeCreate == nil {
		return res, errors.BadRequest.New("DatetimeCreate is required")
	}
	datetimeCreate := p.DatetimeCreate.AsTime()

	// Convert budget
	var budget *CreateAccountBudgetReq
	if p.Budget != nil {
		budget = &CreateAccountBudgetReq{
			Amount:         decimal.NewFromFloat(p.Budget.Amount),
			FixedSum:       decimal.NewFromFloat(p.Budget.FixedSum),
			DaysOffset:     p.Budget.DaysOffset,
			GradualFilling: &p.Budget.GradualFilling,
		}
	}

	var parentAccountID *uuid.UUID
	if len(p.ParentAccountID) != 0 {
		_parentAccountID, err := uuid.FromBytes(p.ParentAccountID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		parentAccountID = &_parentAccountID
	}

	return CreateAccountReq{
		ID:                 id,
		Name:               p.Name,
		IconID:             iconID,
		Type:               accountType,
		Currency:           p.Currency,
		AccountGroupID:     accountGroupID,
		AccountingInHeader: p.AccountingInHeader,
		AccountingInCharts: p.AccountingInCharts,
		DatetimeCreate:     datetimeCreate,
		Remainder:          decimal.NewFromFloat(*p.Remainder),
		Budget:             budget,
		IsParent:           p.IsParent,
		ParentAccountID:    parentAccountID,
	}, nil
}
