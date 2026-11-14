package service

import (
	"github.com/shopspring/decimal"

	"server/internal/utils/errors"
)

// convertPrecision — знаков после запятой в пересчитанных суммах: хватает для криптовалют
const convertPrecision = 8

// convert пересчитывает сумму из валюты счета в валюту ответа по текущему курсу. Курс в базе —
// сколько единиц валюты дают за 1 USD (RUB ≈ 82), поэтому сумма делится на курс своей валюты и
// умножается на курс целевой, как в приложении: amount * targetRate / accountRate
func (d dataset) convert(amount decimal.Decimal, fromCurrency string) (decimal.Decimal, error) {
	if fromCurrency == d.currency {
		return amount, nil
	}
	fromRate, ok := d.rates[fromCurrency]
	if !ok || fromRate.IsZero() {
		return decimal.Zero, errors.InternalServer.New("exchange rate not found").WithParams("currency", fromCurrency)
	}
	return amount.Mul(d.rates[d.currency]).Div(fromRate).Round(convertPrecision), nil
}
