package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/pendingLinkedTransfer/model"
	repoModel "server/internal/modules/pendingLinkedTransfer/repository/model"
)

// DeletePendingLinkedTransfer удаляет перенос — например, вместе с удалением транзакции-инициатора
// (без источника перенос теряет смысл, а sourceTransactionID у клиента — реальный локальный FK).
func (s *PendingLinkedTransferService) DeletePendingLinkedTransfer(ctx context.Context, req model.DeletePendingLinkedTransferReq) error {
	ctx, span := tracer.Start(ctx, "DeletePendingLinkedTransfer")
	defer span.End()

	// Получаем перенос для слепка "до" в аудит-логе
	transferBefore, err := slices.FirstWithError(s.pendingLinkedTransferRepository.GetPendingLinkedTransfers(ctx, repoModel.GetPendingLinkedTransfersReq{ //nolint:exhaustruct
		IDs: []uuid.UUID{req.ID},
	}))
	if err != nil {
		return err
	}

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Удаляем перенос
		if err := s.pendingLinkedTransferRepository.DeletePendingLinkedTransfer(ctxTx, req.ID); err != nil {
			return 0, err
		}

		// Фиксируем удаление переноса в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:         auditLogEntity.PendingLinkedTransfer,
			Method:         auditLogMethod.Delete,
			EntityID:       req.ID.String(),
			Before:         transferBefore,
			After:          nil,
			UserID:         req.Necessary.UserID,
			DeviceID:       req.Necessary.DeviceID,
			AccountGroupID: &transferBefore.AccountGroupID,
		})
	})
}
