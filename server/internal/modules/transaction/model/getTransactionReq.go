package model

import (
	"context"

	"github.com/google/uuid"

	"pkg/datetime"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	"server/internal/modules/transaction/model/transactionType"
)

type GetTransactionsReq struct {
	Necessary       necessary.NecessaryUserInformation
	IDs             []uuid.UUID           `json:"-"`                                                                       // Идентификаторы транзакций
	AccountID       *uuid.UUID            `json:"accountID" schema:"accountID" minimum:"1"`                                // Транзакции какого счета нас интересуют
	Type            *transactionType.Type `json:"type" schema:"type" enums:"consumption,income,transfer"`                  // Тип транзакции
	DateFrom        *datetime.Date        `json:"dateFrom" schema:"dateFrom" format:"date" swaggertype:"primitive,string"` // Дата, от которой начинать учитывать транзакции
	DateTo          *datetime.Date        `json:"dateTo" schema:"dateTo" format:"date" swaggertype:"primitive,string"`     // Дата, до которой учитывать транзакции
	Offset          *uint32               `json:"offset" schema:"offset" minimum:"0"`                                      // Смещение относительно начала списка для пагинации
	Limit           *uint32               `json:"limit" schema:"limit" minimum:"1"`                                        // Количество транзакций в списке для пагинации
	AccountGroupIDs []uuid.UUID           // Идентификаторы групп счетов
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
