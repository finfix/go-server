package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/utils/necessary"

	"server/internal/modules/tag/model"
)

// @Summary Удаление транзакции
// @Description Удаление данных транзакции и изменение баланса счетов
// @Tags tag
// @Security AuthJWT
// @Param Query query model.DeleteTagReq true "model.DeleteTagReq"
// @Produce json
// @Success 200 "Если удаление транзакции прошло успешно, возвращается пустой ответ"
// @Failure 400,401,403,500 {object} errors.Error
// @Router /tag [delete]
func (s *endpoint) deleteTag(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "deleteTag")
	defer span.End()

	var req model.DeleteTagReq

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
	return nil, s.service.DeleteTag(ctx, req)
}
