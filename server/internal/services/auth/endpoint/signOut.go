package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/utils/necessary"

	"server/internal/services/auth/model"
)

// @Summary Выход пользователя из приложения
// @Tags auth
// @Accept json
// @Produce json
// @Security AuthJWT
// @Success 200 "При успешном выходе возвращается null"
// @Failure 400,404,500 {object} errors.Error
// @Router /auth/signOut [post]
func (s *endpoint) signOut(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "signOut")
	defer span.End()

	var req model.SignOutReq

	// Декодируем запрос
	if err := decoder.Decode(ctx, r, &req, decoder.DecodeJSON); err != nil {
		return nil, err
	}

	// Парсим обязательные параметры
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return nil, err
	}

	// Валидируем запрос
	if err := validator.Validate(ctx, req); err != nil {
		return nil, err
	}

	// Вызываем метод сервиса
	return nil, s.service.SignOut(ctx, req)
}
