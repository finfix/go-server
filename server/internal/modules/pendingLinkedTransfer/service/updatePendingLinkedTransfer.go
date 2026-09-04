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

// UpdatePendingLinkedTransfer меняет статус переноса (Completed/Ignored)
func (s *PendingLinkedTransferService) UpdatePendingLinkedTransfer(ctx context.Context, req model.UpdatePendingLinkedTransferReq) error {
	ctx, span := tracer.Start(ctx, "UpdatePendingLinkedTransfer")
	defer span.End()

	// Получаем перенос для слепка "до" в аудит-логе
	transferBefore, err := slices.FirstWithError(s.pendingLinkedTransferRepository.GetPendingLinkedTransfers(ctx, repoModel.GetPendingLinkedTransfersReq{ //nolint:exhaustruct
		IDs: []uuid.UUID{req.ID},
	}))
	if err != nil {
		return err
	}

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Обновляем статус переноса
		if err := s.pendingLinkedTransferRepository.UpdatePendingLinkedTransfer(ctxTx, req.ID, req.ConvertToRepoReq()); err != nil {
			return 0, err
		}

		// Получаем актуальный перенос из БД для слепка "после" в аудит-логе
		transferAfter, err := slices.FirstWithError(s.pendingLinkedTransferRepository.GetPendingLinkedTransfers(ctxTx, repoModel.GetPendingLinkedTransfersReq{ //nolint:exhaustruct
			IDs: []uuid.UUID{req.ID},
		}))
		if err != nil {
			return 0, err
		}

		// Фиксируем изменение переноса в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:         auditLogEntity.PendingLinkedTransfer,
			Method:         auditLogMethod.Update,
			EntityID:       req.ID.String(),
			Before:         transferBefore,
			After:          transferAfter,
			UserID:         req.Necessary.UserID,
			DeviceID:       req.Necessary.DeviceID,
			AccountGroupID: &transferAfter.AccountGroupID,
		})
	})
}
