package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// portfolioValue — стоимость остатков счетов (в их валютах) в валюте ответа по курсам на момент at
func (c seriesCalculator) portfolioValue(balances map[uuid.UUID]decimal.Decimal, at time.Time) (decimal.Decimal, error) {
	total := decimal.Zero
	for id, amount := range balances {
		value, err := c.data.convert(amount, c.data.accounts[id].Currency, at)
		if err != nil {
			return decimal.Zero, err
		}
		total = total.Add(value)
	}
	return total, nil
}
