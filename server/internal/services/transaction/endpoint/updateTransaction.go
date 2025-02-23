package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/utils/necessary"

	"server/internal/services/transaction/model"
)

// @Summary Редактирование транзакции
// @Description Изменение данных транзакции
// @Tags transaction
// @Security AuthJWT
// @Accept json
// @Param Body body model.UpdateTransactionReq true "model.UpdateTransactionReq"
// @Success 200 "При успешном выполнении запроса возвращает пустой ответ"
// @Failure 400,401,403,404,500 {object} errors.Error
// @Router /transaction [patch]
func (s *endpoint) updateTransaction(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "updateTransaction")
	defer span.End()

	var req model.UpdateTransactionReq

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
	return nil, s.service.UpdateTransaction(ctx, req)
}
