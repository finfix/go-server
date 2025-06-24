package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/utils/necessary"

	"server/internal/services/tag/model"
)

// @Summary Создание подкатегории
// @Description Создание подкатегории
// @Tags tag
// @Security AuthJWT
// @Accept json
// @Param Body body model.CreateTagReq true "model.CreateTagReq"
// @Produce json
// @Success 200 {object} model.CreateTagRes
// @Failure 400,401,403,404,500 {object} errors.Error
// @Router /tag [post]
func (s *endpoint) createTag(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "createTag")
	defer span.End()

	var req model.CreateTagReq

	// Декодируем запрос
	if err := decoder.Decode(ctx, r, &req, decoder.DecodeJSON); err != nil {
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
	id, err := s.service.CreateTag(ctx, req)
	if err != nil {
		return nil, err
	}

	return model.CreateTagRes{ID: id}, nil
}
