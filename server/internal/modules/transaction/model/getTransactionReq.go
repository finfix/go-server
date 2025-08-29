package model

import (
	"context"
	"server/internal/enum/transactionType"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"

	"pkg/datetime"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type GetTransactionsReq struct {
	Necessary       necessary.NecessaryUserInformation
	IDs             []uuid.UUID                      `json:"-"`                                                                       // Идентификаторы транзакций
	AccountID       *uuid.UUID                       `json:"accountID" schema:"accountID" minimum:"1"`                                // Транзакции какого счета нас интересуют
	Type            *transactionType.TransactionType `json:"type" schema:"type" enums:"consumption,income,transfer"`                  // Тип транзакции
	DateFrom        *datetime.Date                   `json:"dateFrom" schema:"dateFrom" format:"date" swaggertype:"primitive,string"` // Дата, от которой начинать учитывать транзакции
	DateTo          *datetime.Date                   `json:"dateTo" schema:"dateTo" format:"date" swaggertype:"primitive,string"`     // Дата, до которой учитывать транзакции
	Offset          *uint32                          `json:"offset" schema:"offset" minimum:"0"`                                      // Смещение относительно начала списка для пагинации
	Limit           *uint32                          `json:"limit" schema:"limit" minimum:"1"`                                        // Количество транзакций в списке для пагинации
	AccountGroupIDs []uuid.UUID                      // Идентификаторы групп счетов
}

func (s GetTransactionsReq) Validate(ctx context.Context) error {
	if err := s.Type.Validate(ctx); err != nil {
		return err
	}
	if s.DateFrom != nil && s.DateTo != nil {
		if s.DateFrom.After(s.DateTo.Time) || s.DateFrom.Equal(s.DateTo.Time) {
			return errors.BadRequest.New("date_from must be less than date_to").WithContextParams(ctx)
		}
	}
	return nil
}

// ProtoGetTransactionsReq wrapper for proto request
type ProtoGetTransactionsReq struct {
	*proto.GetTransactionsRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoGetTransactionsReq) ConvertToModel() (GetTransactionsReq, error) {
	var res GetTransactionsReq

	if p.GetTransactionsRequest == nil {
		return res, errors.BadRequest.New("GetTransactionsRequest is required")
	}

	// Parse optional AccountID
	var accountID *uuid.UUID
	if p.AccountID != nil {
		parsedAccountID, err := uuid.FromBytes(p.AccountID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		accountID = &parsedAccountID
	}

	// Convert account group IDs
	accountGroupIDs := make([]uuid.UUID, 0, len(p.AccountGroupIDs))
	for _, idBytes := range p.AccountGroupIDs {
		id, err := uuid.FromBytes(idBytes)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		accountGroupIDs = append(accountGroupIDs, id)
	}

	// Convert optional transaction type
	var transactionTypePtr *transactionType.TransactionType
	if p.Type != nil {
		transactionTypeConverted, err := transactionType.ProtoTransactionType{TransactionType: *p.Type}.ConvertToModel()
		if err != nil {
			return res, err
		}
		transactionTypePtr = &transactionTypeConverted
	}

	// Convert optional dates
	var dateFrom *datetime.Date
	if p.DateFrom != nil {
		dateFrom = &datetime.Date{Time: p.DateFrom.AsTime()}
	}

	var dateTo *datetime.Date
	if p.DateTo != nil {
		dateTo = &datetime.Date{Time: p.DateTo.AsTime()}
	}

	return GetTransactionsReq{
		AccountID:       accountID,
		AccountGroupIDs: accountGroupIDs,
		Type:            transactionTypePtr,
		DateFrom:        dateFrom,
		DateTo:          dateTo,
		Offset:          p.Offset,
		Limit:           p.Limit,
	}, nil
}
