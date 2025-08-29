package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *TransactionService) CheckAccess(ctx context.Context, userID uuid.UUID, transactionIDs []uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "CheckAccess")
	defer span.End()

	// Проверяем доступ пользователя к группам счетов
	accessedAccountGroups, err := s.userService.GetAccessedAccountGroups(ctx, userID)
	if err != nil {
		return err
	}

	// Проверяем доступ пользователя к транзакциям
	return s.transactionRepository.CheckAccess(ctx, accessedAccountGroups, transactionIDs)
}
