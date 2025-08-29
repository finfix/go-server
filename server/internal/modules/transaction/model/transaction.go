package model

import (
	"server/internal/enum/transactionType"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/finfix/go-server-grpc/proto"

	"pkg/datetime"
)

type Transaction struct {
	ID                 uuid.UUID                       `json:"id" db:"id" minimum:"1"`                                                             // Идентификатор транзакции
	Type               transactionType.TransactionType `json:"type" db:"transaction_type" enums:"consumption,income,transfer"`                     // Тип транзакции
	AmountFrom         decimal.Decimal                 `json:"amountFrom" db:"amount_from" minimum:"1"`                                            // Сумма сделки в первой валюте
	AmountTo           decimal.Decimal                 `json:"amountTo" db:"amount_to" minimum:"1"`                                                // Сумма сделки во второй валюте
	Note               string                          `json:"note" db:"note"`                                                                     // Заметка сделки
	AccountFromID      uuid.UUID                       `json:"accountFromID" db:"account_from_id" minimum:"1"`                                     // Идентификатор счета списания
	AccountToID        uuid.UUID                       `json:"accountToID" db:"account_to_id" minimum:"1"`                                         // Идентификатор счета пополнения
	DateTransaction    datetime.Date                   `json:"dateTransaction" db:"date_transaction" format:"date" swaggertype:"primitive,string"` // Дата транзакции (пользовательские)
	IsExecuted         bool                            `json:"isExecuted" db:"is_executed"`                                                        // Исполнена операция или нет (если нет, сделки как бы не существует)
	AccountingInCharts bool                            `json:"accountingInCharts" db:"accounting_in_charts"`                                       // Учитывается ли транзакция в графиках или нет
	CreatedByUserID    uuid.UUID                       `json:"createdByUserID" db:"created_by_user_id" minimum:"1"`                                // Идентификатор пользователя, создавшего транзакцию
	DatetimeCreate     datetime.Time                   `json:"datetimeCreate" db:"datetime_create" format:"date-time"`                             // Дата и время создания транзакции
	AccountGroupID     uuid.UUID                       `json:"accountGroupID" db:"account_group_id"`                                               // Идентификатор группы счетов
}

// ConvertToProto converts Transaction to proto format
func (t Transaction) ConvertToProto() (*proto.Transaction, error) {
	protoType, err := t.Type.ConvertToProto()
	if err != nil {
		return nil, err
	}

	return &proto.Transaction{
		Id:                 t.ID[:],
		Type:               protoType,
		AmountFrom:         t.AmountFrom.InexactFloat64(),
		AmountTo:           t.AmountTo.InexactFloat64(),
		Note:               t.Note,
		AccountFromID:      t.AccountFromID[:],
		AccountToID:        t.AccountToID[:],
		DateTransaction:    timestamppb.New(t.DateTransaction.Time),
		IsExecuted:         t.IsExecuted,
		AccountingInCharts: t.AccountingInCharts,
		CreatedByUserID:    t.CreatedByUserID[:],
		DatetimeCreate:     timestamppb.New(t.DatetimeCreate.Time),
		AccountGroupID:     t.AccountGroupID[:],
	}, nil
}
