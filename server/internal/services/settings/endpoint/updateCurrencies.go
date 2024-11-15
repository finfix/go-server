package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"

	"server/internal/services/settings/model"
)

// @Summary Обновление курсов валют
// @Tags settings
// @Security AuthJWT
// @Success 200 "При успешном выполнении возвращается пустой ответ"
// @Failure 400,401,403,500 {object} errors.Error
// @Router /settings/updateCurrencies [post]
func (s *endpoint) updateCurrencies(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "updateCurrencies")
	defer span.End()

	var req model.UpdateCurrenciesReq

	// Декодируем запрос
	if err := decoder.Decoder(ctx, r, &req); err != nil {
		return nil, err
	}

	return nil, s.service.UpdateCurrencies(ctx, req)
}
