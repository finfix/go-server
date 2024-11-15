package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"

	"server/internal/services/accountGroup/model"
)

// @Summary Получение списка групп счетов
// @Tags accountGroup
// @Security AuthJWT
// @Param Query query model.GetAccountGroupsReq true "model.GetAccountGroupsReq"
// @Produce json
// @Success 200 {object} []model.AccountGroup "[]model.AccountGroup"
// @Failure 400,401,403,404,500 {object} errors.Error
// @Router /accountGroup [get]
func (s *endpoint) getAccountGroups(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "getAccountGroups")
	defer span.End()

	var req model.GetAccountGroupsReq

	// Декодируем запрос
	if err := decoder.Decoder(ctx, r, &req, decoder.DecodeSchema); err != nil {
		return nil, err
	}

	// Вызываем метод сервиса
	return s.service.GetAccountGroups(ctx, req)
}
