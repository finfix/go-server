package model

import (
	"server/internal/enum/accountType"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"pkg/datetime"
)

type Account struct {
	ID                 uuid.UUID               `json:"id" db:"id"`                                                   // Идентификатор счета
	Remainder          decimal.Decimal         `json:"remainder" db:"remainder"`                                     // Остаток средств на счету
	Name               string                  `json:"name" db:"name"`                                               // Название счета
	IconID             uuid.UUID               `json:"iconID" db:"icon_id"`                                          // Идентификатор иконки
	Type               accountType.AccountType `json:"type" db:"account_type" enums:"regular,expense,debt,earnings"` // Тип счета
	Currency           string                  `json:"currency" db:"currency"`                                       // Валюта счета
	Visible            bool                    `json:"visible" db:"visible"`                                         // Видимость счета
	AccountGroupID     uuid.UUID               `json:"accountGroupID" db:"account_group_id"`                         // Идентификатор группы счета
	AccountingInHeader bool                    `json:"accountingInHeader" db:"accounting_in_header"`                 // Будет ли счет учитываться в шапке
	ParentAccountID    *uuid.UUID              `json:"parentAccountID" db:"parent_account_id" validate:"required"`   // Идентификатор родительского аккаунта
	SerialNumber       uint32                  `json:"serialNumber" db:"serial_number"`                              // Порядковый номер счета
	IsParent           bool                    `json:"isParent" db:"is_parent"`                                      // Является ли счет родительским
	CreatedByUserID    uuid.UUID               `json:"createdByUserID" db:"created_by_user_id"`                      // Идентификатор пользователя, создавшего счет
	DatetimeCreate     datetime.Time           `json:"datetimeCreate" db:"datetime_create"`                          // Дата и время создания счета
	AccountingInCharts bool                    `json:"accountingInCharts" db:"accounting_in_charts"`                 // Учитывать ли счет в графиках
	AccountBudget      `json:"budget"`         // Бюджет
}

// ConvertToProto converts internal Account to proto Account
func (a Account) ConvertToProto() (*proto.Account, error) {

	// Convert account type
	protoAccountType, err := a.Type.ConvertToProto()
	if err != nil {
		return nil, err
	}

	// Convert budget
	var protoBudget *proto.AccountBudget
	if true { // Always include budget for now
		protoBudget = &proto.AccountBudget{
			Amount:         a.AccountBudget.Amount.InexactFloat64(),
			DaysOffset:     a.AccountBudget.DaysOffset,
			FixedSum:       a.AccountBudget.FixedSum.InexactFloat64(),
			GradualFilling: a.AccountBudget.GradualFilling,
		}
	}

	var parentAccountID []byte
	if a.ParentAccountID != nil {
		parentAccountID = a.ParentAccountID[:]
	}

	return &proto.Account{
		Id:                 a.ID[:],
		Name:               a.Name,
		Type:               protoAccountType,
		Currency:           a.Currency,
		Remainder:          a.Remainder.InexactFloat64(),
		Visible:            a.Visible,
		AccountingInCharts: a.AccountingInCharts,
		AccountingInHeader: a.AccountingInHeader,
		AccountGroupID:     a.AccountGroupID[:],
		ParentAccountID:    parentAccountID,
		IsParent:           a.IsParent,
		IconID:             a.IconID[:],
		SerialNumber:       a.SerialNumber,
		CreatedByUserID:    a.CreatedByUserID[:],
		DatetimeCreate:     timestamppb.New(a.DatetimeCreate.Time),
		Budget:             protoBudget,
	}, nil
}
