package mcp

import (
	"context"
	"fmt"
	"sort"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shopspring/decimal"

	settingsModel "server/internal/modules/settings/model"
	mcpTransport "server/internal/utils/mcp"
)

// Пределы ответа get_rates
const (
	maxRateCurrencies = 10   // Валют в запросе
	maxRatePoints     = 3000 // Точек на валюту
	defaultRateDays   = 30   // Дней без dateFrom
)

// GetRatesIn — аргументы get_rates
type GetRatesIn struct {
	Currencies []string `json:"currencies" jsonschema:"Коды валют, например RUB, ETH"`
	Base       string   `json:"base,omitempty" jsonschema:"Валюта, за единицу которой указан курс, по умолчанию USD"`
	DateFrom   *string  `json:"dateFrom,omitempty" jsonschema:"Начало периода, yyyy-MM-dd; по умолчанию 30 дней до dateTo"`
	DateTo     *string  `json:"dateTo,omitempty" jsonschema:"Конец периода включительно, yyyy-MM-dd; по умолчанию сегодня"`
}

// GetRatesOut — результат get_rates
type GetRatesOut struct {
	Base   string           `json:"base"`   // За единицу какой валюты курс
	Series []RateSeriesView `json:"series"` // Ряды по валютам
}

// RateSeriesView — история курса одной валюты
type RateSeriesView struct {
	Currency string          `json:"currency"` // Валюта
	Points   []RatePointView `json:"points"`   // Точки по возрастанию даты
}

// RatePointView — точка курса
type RatePointView struct {
	At   string          `json:"at"`   // Момент, RFC 3339
	Rate decimal.Decimal `json:"rate"` // Единиц валюты за 1 base
}

// GetRates история курсов
func (t *SettingsTools) GetRates(ctx context.Context, _ *sdk.CallToolRequest, in GetRatesIn) (*sdk.CallToolResult, GetRatesOut, error) {
	var out GetRatesOut

	// Проверяем аргументы
	if _, err := mcpTransport.Necessary(ctx); err != nil {
		return nil, out, err
	}
	if len(in.Currencies) == 0 || len(in.Currencies) > maxRateCurrencies {
		return nil, out, fmt.Errorf("currencies: от 1 до %d валют", maxRateCurrencies)
	}
	out.Base = in.Base
	if out.Base == "" {
		out.Base = "USD"
	}
	dateTo, err := mcpTransport.ParseOptionalDate("dateTo", in.DateTo)
	if err != nil {
		return nil, out, err
	}
	if dateTo == nil {
		today := time.Now().UTC().Truncate(24 * time.Hour)
		dateTo = &today
	}
	dateFrom, err := mcpTransport.ParseOptionalDate("dateFrom", in.DateFrom)
	if err != nil {
		return nil, out, err
	}
	if dateFrom == nil {
		from := dateTo.AddDate(0, 0, -defaultRateDays)
		dateFrom = &from
	}

	// История валют и базы; база нужна и на моменты до периода, чтобы пересчитать первую точку
	rates, err := t.settingsService.GetCurrencyRates(ctx, settingsModel.GetCurrencyRatesReq{Currencies: append([]string{out.Base}, in.Currencies...), DateTo: dateTo}) //nolint:exhaustruct
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}
	byCurrency := make(map[string][]settingsModel.CurrencyRate)
	for _, rate := range rates {
		byCurrency[rate.Currency] = append(byCurrency[rate.Currency], rate)
	}
	base := byCurrency[out.Base]
	if len(base) == 0 {
		return nil, out, fmt.Errorf("base: нет курсов валюты %s", out.Base)
	}

	// Курс валюты за 1 base на момент точки: делим на действующий курс базы
	for _, currency := range in.Currencies {
		series := RateSeriesView{Currency: currency, Points: []RatePointView{}}
		points := byCurrency[currency]

		// Первая точка — курс, действующий на начало периода, даже если он установлен раньше
		first := sort.Search(len(points), func(i int) bool { return !points[i].ValidAt.Before(*dateFrom) })
		if first > 0 && (first == len(points) || points[first].ValidAt.After(*dateFrom)) {
			first--
		}
		for _, point := range points[first:] {
			if len(series.Points) >= maxRatePoints {
				break
			}
			index := sort.Search(len(base), func(i int) bool { return base[i].ValidAt.After(point.ValidAt) })
			if index == 0 {
				index = 1
			}
			series.Points = append(series.Points, RatePointView{
				At:   point.ValidAt.UTC().Format(time.RFC3339),
				Rate: point.Rate.Div(base[index-1].Rate).Round(8),
			})
		}
		out.Series = append(out.Series, series)
	}
	return nil, out, nil
}
