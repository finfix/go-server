package utils

import "github.com/shopspring/decimal"

func GetRate(rates map[string]decimal.Decimal, currency, currencyRelate string) decimal.Decimal {
	currencyRate := rates[currency]
	currencyRelateRate := rates[currencyRelate]
	if currencyRate.GreaterThan(currencyRelateRate) {
		return currencyRate.Div(currencyRelateRate)

	}
	return currencyRelateRate.Div(currencyRate)
}
