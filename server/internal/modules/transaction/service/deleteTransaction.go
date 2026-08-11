package service

import (
	"context"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/transaction/model"

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

	// Получаем транзакцию для слепка "до" в аудит-логе
	transactionBefore, err := slices.FirstWithError(s.transactionRepository.GetTransactions(ctx, model.GetTransactionsReq{ //nolint:exhaustruct
		IDs: []uuid.UUID{id.ID},
	}))
	if err != nil {
		return err
	}

	return s.generalRepository.WithSyncGate(ctx, id.Necessary.UserID, id.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Удаляем транзакцию
		if err := s.transactionRepository.DeleteTransaction(ctxTx, id.ID, id.Necessary.UserID); err != nil {
			return 0, err
		}

		// Фиксируем удаление транзакции в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:         auditLogEntity.Transaction,
			Method:         auditLogMethod.Delete,
			EntityID:       id.ID.String(),
			Before:         transactionBefore,
			After:          nil,
			UserID:         id.Necessary.UserID,
			DeviceID:       id.Necessary.DeviceID,
			AccountGroupID: &transactionBefore.AccountGroupID,
		})
	})
}
