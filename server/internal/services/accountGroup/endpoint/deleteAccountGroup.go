package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/utils/necessary"

	"server/internal/services/accountGroup/model"
)

// @Summary Удаление группы счетов
// @Tags accountGroup
// @Security AuthJWT
// @Param Query query model.DeleteAccountGroupReq true "model.DeleteAccountGroupReq"
// @Produce json
// @Success 200 "Если удаление группы счетов прошло успешно, возвращается пустой ответ"
// @Failure 400,401,403,404,500 {object} errors.Error
// @Router /accountGroup [delete]
func (s *endpoint) deleteAccountGroup(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "deleteAccountGroup")
	defer span.End()

	var req model.DeleteAccountGroupReq

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
	return nil, s.service.DeleteAccountGroup(ctx, req)
}
