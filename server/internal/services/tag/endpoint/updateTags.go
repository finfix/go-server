package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"

	"server/internal/services/tag/model"
)

// @Summary Редактирование транзакции
// @Description Изменение данных транзакции
// @Tags tag
// @Security AuthJWT
// @Accept json
// @Param Body body model.UpdateTagReq true "model.UpdateTagReq"
// @Success 200 "При успешном выполнении запроса возвращает пустой ответ"
// @Failure 400,401,403,404,500 {object} errors.Error
// @Router /tag [patch]
func (s *endpoint) updateTag(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "updateTag")
	defer span.End()

	var req model.UpdateTagReq

	// Декодируем запрос
	if err := decoder.Decoder(ctx, r, &req, decoder.DecodeJSON); err != nil {
		return nil, err
	}

	// Вызываем метод сервиса
	return nil, s.service.UpdateTag(ctx, req)
}
