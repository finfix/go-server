package network

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/shopspring/decimal"

	"server/internal/utils/errors"
)

// requestCurrencyRates запрашивает курсы у провайдера: курсы — единиц валюты за 1 USD, момент —
// на который провайдер их посчитал
func requestCurrencyRates(ctx context.Context, urlString string, urlValues url.Values) (map[string]decimal.Decimal, time.Time, error) {
	var providerModel struct {
		Meta struct {
			LastUpdatedAt time.Time `json:"last_updated_at"`
		} `json:"meta"`
		Rates map[string]struct {
			Rate decimal.Decimal `json:"value"`
		} `json:"data"`
	}

	// Собираем адрес с параметрами
	uri, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, time.Time{}, errors.InternalServer.Wrap(err).WithContextParams(ctx)
	}
	uri.RawQuery = urlValues.Encode()

	// Отправляем запрос
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, time.Time{}, errors.InternalServer.Wrap(err).WithContextParams(ctx)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, time.Time{}, errors.BadGateway.Wrap(err).WithContextParams(ctx)
	}
	defer resp.Body.Close()

	// Смотрим код ответа и декодируем
	if resp.StatusCode != http.StatusOK {
		return nil, time.Time{}, errors.BadGateway.New("Error while getting currency rates").
			WithContextParams(ctx).
			WithParams("HTTPCode", resp.StatusCode, "url", urlString)
	}
	if err = json.NewDecoder(resp.Body).Decode(&providerModel); err != nil {
		return nil, time.Time{}, errors.InternalServer.Wrap(err).WithContextParams(ctx)
	}

	// Конвертируем полученные данные в нужный формат
	rates := make(map[string]decimal.Decimal, len(providerModel.Rates))
	for currency, rate := range providerModel.Rates {
		rates[currency] = rate.Rate
	}
	return rates, providerModel.Meta.LastUpdatedAt, nil
}
