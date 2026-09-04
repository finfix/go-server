package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/enum/pendingLinkedTransferStatus"
	"server/internal/modules/pendingLinkedTransfer/repository/model"
	"server/internal/modules/pendingLinkedTransfer/repository/pendingLinkedTransferDDL"
)

// CreatePendingLinkedTransfer создаёт требование довнесения транзакции через счёт-мост
func (r *PendingLinkedTransferRepository) CreatePendingLinkedTransfer(ctx context.Context, req model.CreatePendingLinkedTransferReq) error {
	ctx, span := tracer.Start(ctx, "createPendingLinkedTransfer")
	defer span.End()

	return r.db.Exec(ctx, sq.
		Insert(pendingLinkedTransferDDL.Table).
		SetMap(map[string]any{
			pendingLinkedTransferDDL.ColumnID:                  req.ID,
			pendingLinkedTransferDDL.ColumnStatus:              pendingLinkedTransferStatus.Pending,
			pendingLinkedTransferDDL.ColumnSourceTransactionID: req.SourceTransactionID,
			pendingLinkedTransferDDL.ColumnSourceAccountID:     req.SourceAccountID,
			pendingLinkedTransferDDL.ColumnTargetAccountID:     req.TargetAccountID,
			pendingLinkedTransferDDL.ColumnAccountGroupID:      req.AccountGroupID,
		}))
}
