package network

import (
	"context"
	"net/url"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/shopspring/decimal"
)

var tracer = otel.Tracer("/server/internal/modules/settings/network")

// GetCurrencyRates возвращает текущие курсы провайдера (единиц валюты за 1 USD) и момент, на
// который они посчитаны
func GetCurrencyRates(ctx context.Context, apiKey string) (map[string]decimal.Decimal, time.Time, error) {
	ctx, span := tracer.Start(ctx, "GetCurrencyRates")
	defer span.End()

	return requestCurrencyRates(ctx, "https://api.currencyapi.com/v3/latest", url.Values{"apikey": {apiKey}})
}
