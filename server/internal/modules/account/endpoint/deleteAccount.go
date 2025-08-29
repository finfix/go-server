package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/utils/necessary"

	"server/internal/modules/account/model"
)

// @Summary Удаление счета
// @Description Удаление данных по счету
// @Tags account
// @Security AuthJWT
// @Param Query query model.DeleteAccountReq true "model.DeleteAccountReq"
// @Produce json
// @Success 200 "Если удаление счета прошло успешно, возвращается пустой ответ"
// @Failure 400,401,403,404,500 {object} errors.Error
// @Router /account [delete]
func (s *endpoint) deleteAccount(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "deleteAccount")
	defer span.End()

	var req model.DeleteAccountReq

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
	return nil, s.service.DeleteAccount(ctx, req)
}
