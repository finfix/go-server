package service

import (
	"context"

	settingsModel "server/internal/modules/settings/model"
)

// Пределы страницы истории курсов для приложения
const (
	defaultRatesPageSize = 10000 // Точек по умолчанию
	maxRatesPageSize     = 10000 // Наибольшая страница
)

// GetCurrencyRatesPage возвращает страницу истории курсов после момента since. Страница не режет
// момент пополам: точки последнего момента либо все в странице, либо все в следующей — иначе
// продолжение «после последнего момента» потеряло бы часть валют
func (s *SettingsService) GetCurrencyRatesPage(ctx context.Context, req settingsModel.GetCurrencyRatesReq) (res settingsModel.GetCurrencyRatesRes, err error) {
	ctx, span := tracer.Start(ctx, "GetCurrencyRatesPage")
	defer span.End()

	// Размер страницы; берем на одну точку больше, чтобы узнать, есть ли продолжение
	pageSize := req.Limit
	if pageSize <= 0 || pageSize > maxRatesPageSize {
		pageSize = defaultRatesPageSize
	}
	req.Limit = pageSize + 1
	if res.Rates, err = s.settingsRepository.GetCurrencyRates(ctx, req); err != nil {
		return res, err
	}
	if len(res.Rates) <= pageSize {
		return res, nil
	}

	// Обрезаем лишнюю точку и неполный последний момент
	res.HasMore = true
	res.Rates = res.Rates[:pageSize]
	last := res.Rates[len(res.Rates)-1].ValidAt
	cut := len(res.Rates)
	for cut > 0 && res.Rates[cut-1].ValidAt.Equal(last) {
		cut--
	}
	if cut > 0 {
		res.Rates = res.Rates[:cut]
	}
	return res, nil
}
