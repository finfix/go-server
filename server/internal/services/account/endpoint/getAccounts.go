package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/services/account/model"
	"server/internal/utils/necessary"
)

// @Summary Получение счетов по фильтрам
// @Tags account
// @Security AuthJWT
// @Param Query query model.GetAccountsReq false "model.GetAccountsReq"
// @Produce json
// @Success 200 {object} []model.Account
// @Failure 400,401,403,404,500 {object} errors.Error
// @Router /account [get]
func (s *endpoint) getAccounts(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "getAccounts")
	defer span.End()

	var req model.GetAccountsReq

	// Декодируем запрос
	if err := decoder.Decode(ctx, r, &req, decoder.DecodeSchema); err != nil {
		return nil, err
	}

	// Парсим обязательные параметры
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return nil, err
	}

	// Валидируем запрос
	if err := validator.Validate(req); err != nil {
		return nil, err
	}

	// Вызываем метод сервиса
	return s.service.GetAccounts(ctx, req)
}
