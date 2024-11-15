package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"

	"server/internal/services/transaction/model"
)

// @Summary Удаление транзакции
// @Description Удаление данных транзакции и изменение баланса счетов
// @Tags transaction
// @Security AuthJWT
// @Param Query query model.DeleteTransactionReq true "model.DeleteTransactionReq"
// @Produce json
// @Success 200 "Если удаление транзакции прошло успешно, возвращается пустой ответ"
// @Failure 400,401,403,500 {object} errors.Error
// @Router /transaction [delete]
func (s *endpoint) deleteTransaction(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "deleteTransaction")
	defer span.End()

	var req model.DeleteTransactionReq

	// Декодируем запрос
	if err := decoder.Decoder(ctx, r, &req, decoder.DecodeSchema); err != nil {
		return nil, err
	}

	// Вызываем метод сервиса
	return nil, s.service.DeleteTransaction(ctx, req)
}
