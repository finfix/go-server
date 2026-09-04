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

// CreatePendingLinkedTransfer создаёт требование довнесения транзакции через счёт-мост
func (s *PendingLinkedTransferService) CreatePendingLinkedTransfer(ctx context.Context, req model.CreatePendingLinkedTransferReq) error {
	ctx, span := tracer.Start(ctx, "CreatePendingLinkedTransfer")
	defer span.End()

	return s.transactor.WithSyncGate(ctx, req.Necessary.UserID, req.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Создаём перенос
		if err := s.pendingLinkedTransferRepository.CreatePendingLinkedTransfer(ctxTx, req.ConvertToRepoReq()); err != nil {
			return 0, err
		}

		// Получаем созданный перенос из БД для слепка "после" в аудит-логе
		transferAfter, err := slices.FirstWithError(s.pendingLinkedTransferRepository.GetPendingLinkedTransfers(ctxTx, repoModel.GetPendingLinkedTransfersReq{ //nolint:exhaustruct
			IDs: []uuid.UUID{req.ID},
		}))
		if err != nil {
			return 0, err
		}

		// Фиксируем создание переноса в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:         auditLogEntity.PendingLinkedTransfer,
			Method:         auditLogMethod.Create,
			EntityID:       req.ID.String(),
			Before:         nil,
			After:          transferAfter,
			UserID:         req.Necessary.UserID,
			DeviceID:       req.Necessary.DeviceID,
			AccountGroupID: &transferAfter.AccountGroupID,
		})
	})
}
