package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/utils/necessary"

	"server/internal/modules/accountGroup/model"
)

// @Summary Редактирование группы счетов
// @Tags accountGroup
// @Security AuthJWT
// @Accept json
// @Param Body body model.UpdateAccountGroupReq true "model.UpdateAccountGroupReq"
// @Produce json
// @Success 200 "Если редактирование группы счетов прошло успешно, возвращается пустой ответ"
// @Failure 400,401,403,404,500 {object} errors.Error
// @Router /accountGroup [patch]
func (s *endpoint) updateAccountGroup(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "updateAccountGroup")
	defer span.End()

	var req model.UpdateAccountGroupReq

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
	return nil, s.service.UpdateAccountGroup(ctx, req)
}
