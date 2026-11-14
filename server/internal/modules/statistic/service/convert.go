package service

import (
	"time"

	"github.com/shopspring/decimal"

	"server/internal/utils/errors"
)

// convertPrecision — знаков после запятой в пересчитанных суммах: хватает для криптовалют
const convertPrecision = 8

// convert пересчитывает сумму из валюты счета в валюту ответа по курсам на момент at. Курс —
// сколько единиц валюты дают за 1 USD, поэтому сумма делится на курс своей валюты и умножается на
// курс целевой, как в приложении: amount * targetRate / accountRate
func (d dataset) convert(amount decimal.Decimal, fromCurrency string, at time.Time) (decimal.Decimal, error) {
	if fromCurrency == d.currency || amount.IsZero() {
		return amount, nil
	}
	fromRate, ok := d.rateAt(fromCurrency, at)
	if !ok {
		return decimal.Zero, errors.InternalServer.New("exchange rate not found").WithParams("currency", fromCurrency)
	}
	toRate, ok := d.rateAt(d.currency, at)
	if !ok {
		return decimal.Zero, errors.InternalServer.New("exchange rate not found").WithParams("currency", d.currency)
	}
	return amount.Mul(toRate).Div(fromRate).Round(convertPrecision), nil
}
