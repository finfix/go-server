package service

import (
	"time"

	"github.com/shopspring/decimal"
)

// ratePoint — точка истории курса: с момента at валюта стоит rate единиц за 1 USD
type ratePoint struct {
	at   time.Time       // С какого момента действует
	rate decimal.Decimal // Единиц валюты за 1 USD
}
