package model

import (
	"context"
	"server/internal/enum/transactionType"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/finfix/go-server-grpc/proto"

	"pkg/datetime"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	"server/internal/modules/transaction/repository/model"
)

type CreateTransactionReq struct {
	Necessary          necessary.NecessaryUserInformation
	ID uuid.UUID `json:"id" validate:"required"` // Идентификатор транзакции

	Type               transactionType.TransactionType `json:"type" validate:"required"`                                                         // Тип транзакции
	AmountFrom         decimal.Decimal                 `json:"amountFrom" validate:"required" minimum:"1"`                                       // Сумма списания с первого счета
	AmountTo           decimal.Decimal                 `json:"amountTo" validate:"required" minimum:"1"`                                         // Сумма пополнения второго счета (в случаях меж валютной транзакции цифры отличаются)
	Note               string                          `json:"note"`                                                                             // Заметка для транзакции
	AccountFromID      uuid.UUID                       `json:"accountFromID" validate:"required" minimum:"1"`                                    // Идентификатор счета списания
	AccountToID        uuid.UUID                       `json:"accountToID" validate:"required" minimum:"1"`                                      // Идентификатор счета пополнения
	DateTransaction    datetime.Date                   `json:"dateTransaction" validate:"required" format:"date" swaggertype:"primitive,string"` // Дата транзакции
	IsExecuted         *bool                           `json:"isExecuted" validate:"required"`                                                   // Исполнена операция или нет (если нет, сделки как бы не существует)
	TagIDs             []uuid.UUID                     `json:"tagIDs"`                                                                           // Идентификаторы тегов
	DatetimeCreate     datetime.Time                   `json:"datetimeCreate" validate:"required"`                                               // Дата создания транзакции
	AccountingInCharts *bool                           `json:"accountingInCharts" validate:"required"`                                           // Учитывается ли транзакция в графиках или нет
	AccountGroupID     uuid.UUID                       `json:"accountGroupID" validate:"required" minimum:"1"`                                   // Идентификатор группы счетов
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
		ID:                 s.ID,
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

// ProtoCreateTransactionReq wrapper for proto request
type ProtoCreateTransactionReq struct {
	*proto.CreateTransactionRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoCreateTransactionReq) ConvertToModel() (CreateTransactionReq, error) {
	var res CreateTransactionReq

	if p.CreateTransactionRequest == nil {
		return res, errors.BadRequest.New("CreateTransactionRequest is required")
	}

	// Parse ID from bytes
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Parse AccountFromID
	accountFromID, err := uuid.FromBytes(p.AccountFromID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Parse AccountToID
	accountToID, err := uuid.FromBytes(p.AccountToID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Parse AccountGroupID
	accountGroupID, err := uuid.FromBytes(p.AccountGroupID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Convert transaction type
	transactionType, err := transactionType.ProtoTransactionType{TransactionType: p.Type}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Convert datetime
	if p.DatetimeCreate == nil {
		return res, errors.BadRequest.New("DatetimeCreate is required")
	}
	datetimeCreate := datetime.Time{Time: p.DatetimeCreate.AsTime()}

	// Convert date transaction
	if p.DateTransaction == nil {
		return res, errors.BadRequest.New("DateTransaction is required")
	}
	dateTransaction := datetime.Date{Time: p.DateTransaction.AsTime()}

	// Convert tag IDs
	var tagIDs []uuid.UUID
	for _, tagIDBytes := range p.TagIDs {
		tagID, err := uuid.FromBytes(tagIDBytes)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		tagIDs = append(tagIDs, tagID)
	}

	return CreateTransactionReq{
		ID:                 id,
		Type:               transactionType,
		AmountFrom:         decimal.NewFromFloat(p.AmountFrom),
		AmountTo:           decimal.NewFromFloat(p.AmountTo),
		Note:               p.Note,
		AccountFromID:      accountFromID,
		AccountToID:        accountToID,
		DateTransaction:    dateTransaction,
		IsExecuted:         &p.IsExecuted,
		TagIDs:             tagIDs,
		DatetimeCreate:     datetimeCreate,
		AccountingInCharts: &p.AccountingInCharts,
		AccountGroupID:     accountGroupID,
	}, nil
}
