package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"pkg/datetime"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	"server/internal/modules/transaction/model/transactionType"
	"server/internal/modules/transaction/repository/model"
)

type CreateTransactionReq struct {
	Necessary          necessary.NecessaryUserInformation
	Type               transactionType.Type `json:"type" validate:"required"`                                                         // Тип транзакции
	AmountFrom         decimal.Decimal      `json:"amountFrom" validate:"required" minimum:"1"`                                       // Сумма списания с первого счета
	AmountTo           decimal.Decimal      `json:"amountTo" validate:"required" minimum:"1"`                                         // Сумма пополнения второго счета (в случаях меж валютной транзакции цифры отличаются)
	Note               string               `json:"note"`                                                                             // Заметка для транзакции
	AccountFromID      uuid.UUID            `json:"accountFromID" validate:"required" minimum:"1"`                                    // Идентификатор счета списания
	AccountToID        uuid.UUID            `json:"accountToID" validate:"required" minimum:"1"`                                      // Идентификатор счета пополнения
	DateTransaction    datetime.Date        `json:"dateTransaction" validate:"required" format:"date" swaggertype:"primitive,string"` // Дата транзакции
	IsExecuted         *bool                `json:"isExecuted" validate:"required"`                                                   // Исполнена операция или нет (если нет, сделки как бы не существует)
	TagIDs             []uuid.UUID          `json:"tagIDs"`                                                                           // Идентификаторы тегов
	DatetimeCreate     datetime.Time        `json:"datetimeCreate" validate:"required"`                                               // Дата создания транзакции
	AccountingInCharts *bool                `json:"accountingInCharts" validate:"required"`                                           // Учитывается ли транзакция в графиках или нет
	AccountGroupID     uuid.UUID            `json:"accountGroupID" validate:"required" minimum:"1"`                                   // Идентификатор группы счетов
}

func (s CreateTransactionReq) Validate(ctx context.Context) error {
	// Валидируем поля
	if err := s.Type.Validate(ctx); err != nil {
		return err
	}
	if s.AmountFrom.LessThanOrEqual(decimal.Zero) || s.AmountTo.LessThanOrEqual(decimal.Zero) {
		return errors.BadRequest.New("amountFrom and amountTo must be greater than 0").
			WithContextParams(ctx)
	}
	return nil
}

func (s *CreateTransactionReq) ConvertToRepoReq() model.CreateTransactionReq {
	return model.CreateTransactionReq{
		Type:               s.Type,
		AmountFrom:         s.AmountFrom,
		AmountTo:           s.AmountTo,
		Note:               s.Note,
		AccountFromID:      s.AccountFromID,
		AccountToID:        s.AccountToID,
		DateTransaction:    s.DateTransaction,
		IsExecuted:         *s.IsExecuted,
		CreatedByUserID:    s.Necessary.UserID,
		DatetimeCreate:     s.DatetimeCreate.Time,
		AccountingInCharts: *s.AccountingInCharts,
		AccountGroupID:     s.AccountGroupID,
	}
}
