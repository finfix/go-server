package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Series — ряд графика: значение на каждый бакет и сводка
type Series struct {
	Key           string            // Ключ ряда: income, expense, delta, balance или идентификатор счета/тега
	Name          string            // Человеческое название
	AccountID     *uuid.UUID        // Счет ряда, если разбивка по счетам
	TagID         *uuid.UUID        // Тег ряда, если разбивка по тегам
	Values        []decimal.Decimal // Значения по бакетам
	Summary       SeriesSummary     // Сводка по бакетам
	MonthlyBudget *decimal.Decimal  // Действующий месячный бюджет расходного счета на конец периода
}
