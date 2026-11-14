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
	Rate   decimal.Decimal `json:"rate"`   // Текущий курс к базовой валюте курсов
}

// Register добавляет инструменты справочников на MCP-сервер
func Register(server *sdk.Server, settingsService SettingsService) {
	tools := &SettingsTools{settingsService: settingsService}
	mcpTransport.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "get_currencies",
		Description: "Валюты с текущими курсами. Сумма X в валюте A равна X * rate(A) / rate(B) в валюте B. История курсов не хранится",
		Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}, //nolint:exhaustruct
	}, tools.GetCurrencies)
}
