package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	mcpTransport "server/internal/utils/mcp"
)

// GetCurrencies справочник валют
func (t *SettingsTools) GetCurrencies(ctx context.Context, _ *sdk.CallToolRequest, _ GetCurrenciesIn) (*sdk.CallToolResult, GetCurrenciesOut, error) {
	var out GetCurrenciesOut

	// Справочник общий, но вызов — только авторизованному пользователю
	if _, err := mcpTransport.Necessary(ctx); err != nil {
		return nil, out, err
	}

	// Вызываем метод сервиса
	currencies, err := t.settingsService.GetCurrencies(ctx)
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Конвертируем ответ
	out.Currencies = make([]CurrencyView, 0, len(currencies.Currencies))
	for _, currency := range currencies.Currencies {
		out.Currencies = append(out.Currencies, CurrencyView{Code: currency.Slug, Name: currency.Name, Symbol: currency.Symbol, Rate: currency.Rate})
	}
	return nil, out, nil
}
