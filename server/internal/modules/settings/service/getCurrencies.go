package service

import (
	"context"

	settingsModel "server/internal/modules/settings/model"
)

func (s *SettingsService) GetCurrencies(ctx context.Context) ([]settingsModel.Currency, error) {
	ctx, span := tracer.Start(ctx, "GetCurrencies")
	defer span.End()

	return s.settingsRepository.GetCurrencies(ctx)
}
