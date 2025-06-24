package endpoint

import (
	"context"
	"net/http"

	_ "server/internal/services/user/model" //nolint:golint
)

// @Summary Получение списка валют
// @Tags settings
// @Security AuthJWT
// @Produce json
// @Success 200 {object} []model.Currency
// @Failure 401,500 {object} errors.Error
// @Router /settings/currencies [get]
func (s *endpoint) getCurrencies(ctx context.Context, _ *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "getCurrencies")
	defer span.End()

	// Вызываем метод сервиса
	return s.service.GetCurrencies(ctx)
}
