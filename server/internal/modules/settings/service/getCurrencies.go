package service

import (
	"context"

	settingsModel "server/internal/modules/settings/model"
)

func (s *SettingsService) GetCurrencies(ctx context.Context) (res settingsModel.GetCurrenciesRes, err error) {
	ctx, span := tracer.Start(ctx, "GetCurrencies")
	defer span.End()

	currencies, err := s.settingsRepository.GetCurrencies(ctx)
	if err != nil {
		return res, err
	}

	return settingsModel.GetCurrenciesRes{Currencies: currencies}, nil
}
