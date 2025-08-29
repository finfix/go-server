package service

import (
	"context"

	"server/internal/services/transaction/model"

	"github.com/google/uuid"
)

// DeleteTransaction удаляет транзакцию
func (s *TransactionService) DeleteTransaction(ctx context.Context, id model.DeleteTransactionReq) error {
	ctx, span := tracer.Start(ctx, "DeleteTransaction")
	defer span.End()

	// Проверяем доступ пользователя к транзакции
	if err := s.CheckAccess(ctx, id.Necessary.UserID, []uuid.UUID{id.ID}); err != nil {
		return err
	}

	// Удаляем транзакцию
	return s.transactionRepository.DeleteTransaction(ctx, id.ID, id.Necessary.UserID)
}
