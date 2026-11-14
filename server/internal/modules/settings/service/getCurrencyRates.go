package service

import (
	"context"

	settingsModel "server/internal/modules/settings/model"
)

// GetCurrencyRates возвращает историю курсов: каждая точка действует до следующей точки той же валюты
func (s *SettingsService) GetCurrencyRates(ctx context.Context, req settingsModel.GetCurrencyRatesReq) ([]settingsModel.CurrencyRate, error) {
	ctx, span := tracer.Start(ctx, "GetCurrencyRates")
	defer span.End()

	return s.settingsRepository.GetCurrencyRates(ctx, req)
}
