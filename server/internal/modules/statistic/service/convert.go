package service

import (
	"github.com/shopspring/decimal"

	"server/internal/utils/errors"
)

// convertPrecision — знаков после запятой в пересчитанных суммах: хватает для криптовалют
const convertPrecision = 8

// convert пересчитывает сумму из валюты счета в валюту ответа по текущему курсу, как
// pkg/currencyConverter: через базовую валюту курсов
func (d dataset) convert(amount decimal.Decimal, fromCurrency string) (decimal.Decimal, error) {
	if fromCurrency == d.currency {
		return amount, nil
	}
	fromRate, ok := d.rates[fromCurrency]
	if !ok || fromRate.IsZero() {
		return decimal.Zero, errors.InternalServer.New("exchange rate not found").WithParams("currency", fromCurrency)
	}
	return amount.Mul(fromRate).Div(d.rates[d.currency]).Round(convertPrecision), nil
}
