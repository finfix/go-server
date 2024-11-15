package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"

	"server/internal/services/account/model"
)

// @Summary Создание счета
// @Description Создается новый счет, если у него есть остаток, то создается транзакция от нулевого счета для баланса
// @Tags account
// @Security AuthJWT
// @Accept json
// @Param Body body model.CreateAccountReq true "model.CreateAccountReq"
// @Produce json
// @Success 200 {object} model.CreateAccountRes
// @Failure 400,401,403,500 {object} errors.Error
// @Router /account [post]
func (s *endpoint) createAccount(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "createAccount")
	defer span.End()

	var req model.CreateAccountReq

	// Декодируем запрос
	if err := decoder.Decoder(ctx, r, &req, decoder.DecodeJSON); err != nil {
		return nil, err
	}

	// Вызываем метод сервиса
	return s.service.CreateAccount(ctx, req)
}
