package network

import (
	"context"
	"net/url"
	"time"

	"github.com/shopspring/decimal"
)

// GetHistoricalCurrencyRates возвращает курсы провайдера на конец дня date (UTC) и момент, на
// который они посчитаны
func GetHistoricalCurrencyRates(ctx context.Context, apiKey string, date time.Time) (map[string]decimal.Decimal, time.Time, error) {
	ctx, span := tracer.Start(ctx, "GetHistoricalCurrencyRates")
	defer span.End()

	return requestCurrencyRates(ctx, "https://api.currencyapi.com/v3/historical", url.Values{"apikey": {apiKey}, "date": {date.UTC().Format("2006-01-02")}})
}
