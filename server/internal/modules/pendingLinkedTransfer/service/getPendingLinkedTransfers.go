package service

import (
	"context"

	"server/internal/modules/pendingLinkedTransfer/model"
)

// GetPendingLinkedTransfers возвращает переносы по фильтрам
func (s *PendingLinkedTransferService) GetPendingLinkedTransfers(ctx context.Context, req model.GetPendingLinkedTransfersReq) ([]model.PendingLinkedTransfer, error) {
	ctx, span := tracer.Start(ctx, "GetPendingLinkedTransfers")
	defer span.End()

	return s.pendingLinkedTransferRepository.GetPendingLinkedTransfers(ctx, req.ConvertToRepoReq())
}
