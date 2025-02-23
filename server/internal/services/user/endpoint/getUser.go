package endpoint

import (
	"context"
	"net/http"

	"pkg/errors"
	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/utils/necessary"

	"server/internal/services/user/model"
)

// @Summary Получение данных пользователя
// @Tags user
// @Security AuthJWT
// @Produce json
// @Success 200 {object} model.User
// @Failure 401,404,500 {object} errors.Error
// @Router /user [get]
func (s *endpoint) getUser(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "getUser")
	defer span.End()

	var req model.GetUsersReq

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
	users, err := s.service.GetUsers(ctx, req)
	if err != nil {
		return nil, errors.InternalServer.Wrap(ctx, err)
	}

	if len(users) == 0 {
		return nil, errors.InternalServer.New(ctx, "Пользователь не найден",
			errors.ParamsOption("UserID", req.Necessary.UserID),
		)
	}

	// Конвертируем ответ во внутреннюю структуру
	return users[0], nil
}
