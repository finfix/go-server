package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/utils/necessary"

	"server/internal/modules/tag/model"
)

// @Summary Получение всех транзакций
// @Description Получение всех транзакций по фильтрам
// @Tags tag
// @Security AuthJWT
// @Param Query query model.GetTagsReq true "model.CreateTagReq"
// @Produce json
// @Success 200 {object} []model.Tag
// @Failure 400,404,500 {object} errors.Error
// @Router /tag [get]
func (s *endpoint) getTags(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "getTags")
	defer span.End()

	var req model.GetTagsReq

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
	return s.service.GetTags(ctx, req)
}
