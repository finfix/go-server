package model

import (
	"context"

	"pkg/pointer"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"pkg/datetime"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type UpdateTransactionReq struct {
	Necessary          necessary.NecessaryUserInformation
	ID                 uuid.UUID        `json:"id" validate:"required" minimum:"1"`                           // Идентификатор транзакции
	AmountFrom         *decimal.Decimal `json:"amountFrom" minimum:"1"`                                       // Сумма списания с первого счета
	AmountTo           *decimal.Decimal `json:"amountTo" minimum:"1"`                                         // Сумма пополнения второго счета
	Note               *string          `json:"note"`                                                         // Заметка для транзакции
	AccountFromID      *uuid.UUID       `json:"accountFromID" minimum:"1"`                                    // Идентификатор счета списания
	AccountToID        *uuid.UUID       `json:"accountToID" minimum:"1"`                                      // Идентификатор счета пополнения
	DateTransaction    *datetime.Date   `json:"dateTransaction" format:"date" swaggertype:"primitive,string"` // Дата транзакции
	IsExecuted         *bool            `json:"isExecuted"`                                                   // Исполнена операция или нет (если нет, сделки как бы не существует)
	TagIDs             *[]uuid.UUID     `json:"tagIDs"`                                                       // Идентификаторы тегов
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

// ProtoUpdateTransactionReq wrapper for proto request
type ProtoUpdateTransactionReq struct {
	*proto.UpdateTransactionRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoUpdateTransactionReq) ConvertToModel() (UpdateTransactionReq, error) {
	var res UpdateTransactionReq

	if p.UpdateTransactionRequest == nil {
		return res, errors.BadRequest.New("UpdateTransactionRequest is required")
	}

	// Parse ID from bytes
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Parse optional AccountFromID
	var accountFromID *uuid.UUID
	if p.AccountFromID != nil {
		parsedAccountFromID, err := uuid.FromBytes(p.AccountFromID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		accountFromID = &parsedAccountFromID
	}

	// Parse optional AccountToID
	var accountToID *uuid.UUID
	if p.AccountToID != nil {
		parsedAccountToID, err := uuid.FromBytes(p.AccountToID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		accountToID = &parsedAccountToID
	}

	// Convert optional amounts
	var amountFrom *decimal.Decimal
	if p.AmountFrom != nil {
		amountFrom = pointer.Pointer(decimal.NewFromFloat(*p.AmountFrom))
	}

	var amountTo *decimal.Decimal
	if p.AmountTo != nil {
		amountTo = pointer.Pointer(decimal.NewFromFloat(*p.AmountTo))
	}

	// Convert optional date transaction
	var dateTransaction *datetime.Date
	if p.DateTransaction != nil {
		dateTransaction = &datetime.Date{Time: p.DateTransaction.AsTime()}
	}

	// Convert optional tag IDs
	var tagIDs *[]uuid.UUID
	if len(p.TagIDs) > 0 {
		convertedTagIDs := make([]uuid.UUID, 0, len(p.TagIDs))
		for _, tagIDBytes := range p.TagIDs {
			tagID, err := uuid.FromBytes(tagIDBytes)
			if err != nil {
				return res, errors.BadRequest.Wrap(err)
			}
			convertedTagIDs = append(convertedTagIDs, tagID)
		}
		tagIDs = &convertedTagIDs
	}

	return UpdateTransactionReq{
		ID:                 id,
		AmountFrom:         amountFrom,
		AmountTo:           amountTo,
		Note:               p.Note,
		AccountFromID:      accountFromID,
		AccountToID:        accountToID,
		DateTransaction:    dateTransaction,
		IsExecuted:         p.IsExecuted,
		TagIDs:             tagIDs,
		AccountingInCharts: p.AccountingInCharts,
	}, nil
}
