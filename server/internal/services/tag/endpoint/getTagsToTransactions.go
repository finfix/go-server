package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/utils/necessary"

	"server/internal/services/tag/model"
)

// @Summary Получение всех связей между подкатегориями и транзакциями
// @Tags tag
// @Security AuthJWT
// @Accept json
// @Param Body body model.GetTagsToTransactionsReq true "model.GetTagsToTransactionsReq"
// @Success 200 "При успешном выполнении запроса возвращает пустой ответ"
// @Failure 400,401,403,404,500 {object} errors.Error
// @Router /tag/to_transactions [get]
func (s *endpoint) getTagsToTransaction(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "getTagsToTransaction")
	defer span.End()

	var req model.GetTagsToTransactionsReq

	// Декодируем запрос
	if err := decoder.Decode(ctx, r, &req, decoder.DecodeSchema); err != nil {
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
	return s.service.GetTagsToTransactions(ctx, req)
}
