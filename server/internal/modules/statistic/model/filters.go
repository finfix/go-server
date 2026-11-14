package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"server/internal/enum/transactionType"
)

// Filters — фильтры транзакций, как на экране фильтров приложения, плюс фильтр по сумме
type Filters struct {
	SearchText         string                            // Подстрока заметки без учета регистра
	DateFrom           *time.Time                        // Дата начала включительно
	DateTo             *time.Time                        // Дата конца включительно
	Types              []transactionType.TransactionType // Типы транзакций
	Currencies         []string                          // Валюты счета списания или пополнения
	AccountIDs         []uuid.UUID                       // Счета (с дочерними), списание или пополнение
	ExcludedAccountIDs []uuid.UUID                       // Исключенные счета (с дочерними)
	TagIDs             []uuid.UUID                       // Теги: транзакция с любым из них
	ExcludedTagIDs     []uuid.UUID                       // Исключенные теги: транзакция без любого из них
	AccountGroupIDs    []uuid.UUID                       // Группы счетов, пусто — все доступные
	AmountMin          *decimal.Decimal                  // Минимальная сумма в валюте счета (любая сторона)
	AmountMax          *decimal.Decimal                  // Максимальная сумма в валюте счета (любая сторона)
}
