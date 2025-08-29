package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/utils/necessary"

	"server/internal/modules/auth/model"
)

// @Summary Обновление токенов
// @Tags auth
// @Accept json
// @Security AuthJWT
// @Param Dody body model.RefreshTokensReq true "model.RefreshTokensReq"
// @Produce json
// @Success 200 {object} model.RefreshTokensRes
// @Failure 400,401,500 {object} errors.Error
// @Router /auth/refreshTokens [post]
func (s *endpoint) refreshTokens(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "refreshTokens")
	defer span.End()

	var req model.RefreshTokensReq

	// Декодируем запрос
	if err := decoder.Decode(ctx, r, &req, decoder.DecodeJSON); err != nil {
		return nil, err
	}

	req.Device.IPAddress = r.Header.Get("X-Real-IP")
	req.Device.UserAgent = r.Header.Get("User-Agent")

	// Парсим обязательные параметры
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return nil, err
	}

	// Валидируем запрос
	if err := validator.Validate(req); err != nil {
		return nil, err
	}

	// Вызываем метод сервиса
	return s.service.RefreshTokens(ctx, req)
}
