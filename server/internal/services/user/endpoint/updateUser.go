package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"

	"server/internal/services/user/model"
)

// @Summary Редактирование пользователя
// @Description Изменение данных пользователя
// @Tags user
// @Security AuthJWT
// @Accept json
// @Param Body body model.UpdateUserReq true "model.UpdateUserReq"
// @Produce json
// @Success 200 "Если редактирование пользователя прошло успешно, возвращается пустой ответ"
// @Failure 400,401,403,404,500 {object} errors.Error
// @Router /user [patch]
func (s *endpoint) updateUser(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "updateUser")
	defer span.End()

	var req model.UpdateUserReq

	// Декодируем запрос
	if err := decoder.Decoder(ctx, r, &req, decoder.DecodeJSON); err != nil {
		return nil, err
	}

	// Вызываем метод сервиса
	return nil, s.service.UpdateUser(ctx, req)
}
