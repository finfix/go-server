package model

import (
	"context"

	"github.com/shopspring/decimal"

	"pkg/datetime"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type UpdateTransactionReq struct {
	Necessary          necessary.NecessaryUserInformation
	ID                 uint32           `json:"id" validate:"required" minimum:"1"`                           // Идентификатор транзакции
	AmountFrom         *decimal.Decimal `json:"amountFrom" minimum:"1"`                                       // Сумма списания с первого счета
	AmountTo           *decimal.Decimal `json:"amountTo" minimum:"1"`                                         // Сумма пополнения второго счета
	Note               *string          `json:"note"`                                                         // Заметка для транзакции
	AccountFromID      *uint32          `json:"accountFromID" minimum:"1"`                                    // Идентификатор счета списания
	AccountToID        *uint32          `json:"accountToID" minimum:"1"`                                      // Идентификатор счета пополнения
	DateTransaction    *datetime.Date   `json:"dateTransaction" format:"date" swaggertype:"primitive,string"` // Дата транзакции
	IsExecuted         *bool            `json:"isExecuted"`                                                   // Исполнена операция или нет (если нет, сделки как бы не существует)
	TagIDs             *[]uint32        `json:"tagIDs"`                                                       // Идентификаторы тегов
	AccountingInCharts *bool            `json:"accountingInCharts"`                                           // Учитывается ли транзакция в графиках или нет
}

func (s UpdateTransactionReq) Validate(ctx context.Context) error {
	if s.AmountFrom != nil && s.AmountFrom.LessThanOrEqual(decimal.Zero) {
		return errors.BadRequest.New("amountFrom must be greater than 0").WithContextParams(ctx)
	}
	if s.AmountTo != nil && s.AmountTo.LessThanOrEqual(decimal.Zero) {
		return errors.BadRequest.New("amountTo must be greater than 0").WithContextParams(ctx)
	}
	return nil
}
