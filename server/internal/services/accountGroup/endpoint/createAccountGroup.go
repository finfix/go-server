package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"

	"server/internal/services/accountGroup/model"
)

// @Summary Создание группы счетов
// @Description Создается новая группа счетов
// @Tags accountGroup
// @Security AuthJWT
// @Accept json
// @Param Body body model.CreateAccountGroupReq true "model.CreateAccountGroupReq"
// @Produce json
// @Success 200 {object} model.CreateAccountGroupRes
// @Failure 400,401,403,500 {object} errors.Error
// @Router /accountGroup [post]
func (s *endpoint) createAccountGroup(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "createAccountGroup")
	defer span.End()

	var req model.CreateAccountGroupReq

	// Декодируем запрос
	if err := decoder.Decoder(ctx, r, &req, decoder.DecodeJSON); err != nil {
		return nil, err
	}

	// Вызываем метод сервиса
	return s.service.CreateAccountGroup(ctx, req)
}
