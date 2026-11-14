// Package mcp — инструмент MCP для справочника валют, только чтение
package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shopspring/decimal"

	settingsModel "server/internal/modules/settings/model"
	settingsService "server/internal/modules/settings/service"
	mcpTransport "server/internal/utils/mcp"
)

var _ SettingsService = new(settingsService.SettingsService)

// SettingsService — интерфейс справочников
type SettingsService interface {
	GetCurrencies(context.Context) (settingsModel.GetCurrenciesRes, error)
	GetCurrencyRates(context.Context, settingsModel.GetCurrencyRatesReq) ([]settingsModel.CurrencyRate, error)
}

// SettingsTools — обработчики инструментов справочников
type SettingsTools struct {
	settingsService SettingsService
}

// GetCurrenciesIn — аргументы get_currencies
type GetCurrenciesIn struct{}

// GetCurrenciesOut — результат get_currencies
type GetCurrenciesOut struct {
	Currencies []CurrencyView `json:"currencies"` // Валюты
}

// CurrencyView — валюта в ответе инструмента
type CurrencyView struct {
	Code   string          `json:"code"`   // Код, например RUB
	Name   string          `json:"name"`   // Название
	Symbol string          `json:"symbol"` // Символ
	Rate   decimal.Decimal `json:"rate"`   // Текущий курс: единиц валюты за 1 USD
}

// Register добавляет инструменты справочников на MCP-сервер
func Register(server *sdk.Server, settingsService SettingsService) {
	tools := &SettingsTools{settingsService: settingsService}
	mcpTransport.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "get_currencies",
		Description: "Валюты с последними курсами. rate — единиц валюты за 1 USD; сумма X в валюте A равна X * rate(B) / rate(A) в валюте B. История курсов не хранится",
		Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}, //nolint:exhaustruct
	}, tools.GetCurrencies)
	mcpTransport.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name: "get_rates",
		Description: "История курсов: по каждой валюте точки (дата, курс) за период. Курс указан в единицах валюты за 1 base (по умолчанию USD): " +
			"USD/RUB 82 — значит 82 рубля за доллар. Точка дня — курс на конец дня по UTC. Без dateFrom — последние 30 дней. Не больше 10 валют и 3 000 точек на валюту",
		Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}, //nolint:exhaustruct
	}, tools.GetRates)
}
