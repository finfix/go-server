package model

import (
	"github.com/shopspring/decimal"
)

type AccountBudget struct {
	Amount         decimal.Decimal `json:"amount" db:"budget_amount"`                  // Сумма бюджета
	FixedSum       decimal.Decimal `json:"fixedSum" db:"budget_fixed_sum"`             // Фиксированная сумма
	DaysOffset     uint32          `json:"daysOffset" db:"budget_days_offset"`         // Смещение в днях
	GradualFilling bool            `json:"gradualFilling" db:"budget_gradual_filling"` // Постепенное пополнение
}
