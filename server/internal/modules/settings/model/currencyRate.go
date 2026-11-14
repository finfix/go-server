package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// CurrencyRate — точка истории курса: с момента ValidAt валюта стоит Rate единиц за 1 USD
type CurrencyRate struct {
	Currency string          `db:"currency"` // Код валюты
	ValidAt  time.Time       `db:"valid_at"` // С какого момента действует курс
	Rate     decimal.Decimal `db:"rate"`     // Единиц валюты за 1 USD
}
