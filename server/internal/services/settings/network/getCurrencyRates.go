package network

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/shopspring/decimal"

	"server/internal/utils/errors"
)

var tracer = otel.Tracer("/server/internal/services/settings/network")

func GetCurrencyRates(ctx context.Context, apiKey string) (map[string]decimal.Decimal, error) {
	ctx, span := tracer.Start(ctx, "GetCurrencyRates")
	defer span.End()

	var providerModel struct {
		Meta struct {
			LastUpdatedAt time.Time `json:"last_updated_at"`
		} `json:"meta"`
		Rates map[string]struct {
			Rate decimal.Decimal `json:"value"`
		} `json:"data"`
	}

	urlValues := url.Values{}

	// URL для получения курсов валют
	urlString := "https://api.currencyapi.com/v3/latest"

	// Параметры запроса
	urlValues.Add("apikey", apiKey)

	uri, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, errors.InternalServer.Wrap(err).WithContextParams(ctx)
	}
	uri.RawQuery = urlValues.Encode()

	// Отправляем запрос
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, errors.InternalServer.Wrap(err).WithContextParams(ctx)
	}

	req = req.WithContext(ctx)
	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.BadGateway.Wrap(err).WithContextParams(ctx)
	}
	defer resp.Body.Close()

	// Смотрим код ответа
	switch resp.StatusCode {
	case http.StatusOK:

		// Декодируем ответ
		if err = json.NewDecoder(resp.Body).Decode(&providerModel); err != nil {
			return nil, errors.InternalServer.Wrap(err).WithContextParams(ctx)
		}
	default:
		return nil, errors.BadGateway.New("Error while getting currency rates").
			WithContextParams(ctx).
			WithParams("HTTPCode", resp.StatusCode)
	}

	rates := make(map[string]decimal.Decimal, len(providerModel.Rates))

	// Конвертируем полученные данные в нужный формат
	for currency, rate := range providerModel.Rates {
		rates[currency] = rate.Rate
	}

	return rates, nil
}
