package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"pkg/datetime"

	"server/internal/services/account/model/accountType"
)

type Account struct {
	ID                 uuid.UUID        `json:"id" db:"id"`                                                     // Идентификатор счета
	Remainder          decimal.Decimal  `json:"remainder" db:"remainder"`                                       // Остаток средств на счету
	Name               string           `json:"name" db:"name"`                                                 // Название счета
	IconID             uuid.UUID        `json:"iconID" db:"icon_id"`                                            // Идентификатор иконки
	Type               accountType.Type `json:"type" db:"type_signatura" enums:"regular,expense,debt,earnings"` // Тип счета
	Currency           string           `json:"currency" db:"currency_signatura"`                               // Валюта счета
	Visible            bool             `json:"visible" db:"visible"`                                           // Видимость счета
	AccountGroupID     uuid.UUID        `json:"accountGroupID" db:"account_group_id"`                           // Идентификатор группы счета
	AccountingInHeader bool             `json:"accountingInHeader" db:"accounting_in_header"`                   // Будет ли счет учитываться в шапке
	ParentAccountID    *uuid.UUID       `json:"parentAccountID" db:"parent_account_id" validate:"required"`     // Идентификатор родительского аккаунта
	SerialNumber       uint32           `json:"serialNumber" db:"serial_number"`                                // Порядковый номер счета
	IsParent           bool             `json:"isParent" db:"is_parent"`                                        // Является ли счет родительским
	CreatedByUserID    uuid.UUID        `json:"createdByUserID" db:"created_by_user_id"`                        // Идентификатор пользователя, создавшего счет
	DatetimeCreate     datetime.Time    `json:"datetimeCreate" db:"datetime_create"`                            // Дата и время создания счета
	AccountingInCharts bool             `json:"accountingInCharts" db:"accounting_in_charts"`                   // Учитывать ли счет в графиках
	AccountBudget      `json:"budget"`  // Бюджет
}
