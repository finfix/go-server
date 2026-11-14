package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"server/internal/enum/transactionType"
)

// TransactionItem — транзакция ленты с именами счетов и тегами
type TransactionItem struct {
	ID              uuid.UUID                       // Идентификатор
	Type            transactionType.TransactionType // Тип
	Date            time.Time                       // Дата
	AccountFromID   uuid.UUID                       // Счет списания
	AccountFromName string                          // Название счета списания
	AmountFrom      decimal.Decimal                 // Сумма списания в валюте счета
	CurrencyFrom    string                          // Валюта счета списания
	AccountToID     uuid.UUID                       // Счет пополнения
	AccountToName   string                          // Название счета пополнения
	AmountTo        decimal.Decimal                 // Сумма пополнения в валюте счета
	CurrencyTo      string                          // Валюта счета пополнения
	Note            string                          // Заметка
	TagIDs          []uuid.UUID                     // Теги
	AccountGroupID  uuid.UUID                       // Группа счетов
	CreatedByUserID uuid.UUID                       // Автор
}
