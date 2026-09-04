package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/pendingLinkedTransfer/model"
	"server/internal/modules/pendingLinkedTransfer/repository/pendingLinkedTransferDDL"
	repoModel "server/internal/modules/pendingLinkedTransfer/repository/model"
	"server/internal/utils/errors"
)

// GetPendingLinkedTransfers возвращает переносы, удовлетворяющие фильтрам
func (r *PendingLinkedTransferRepository) GetPendingLinkedTransfers(ctx context.Context, req repoModel.GetPendingLinkedTransfersReq) (transfers []model.PendingLinkedTransfer, err error) {
	ctx, span := tracer.Start(ctx, "getPendingLinkedTransfers")
	defer span.End()

	filters := make(sq.Eq)

	if len(req.IDs) != 0 {
		filters[pendingLinkedTransferDDL.ColumnID] = req.IDs
	}
	if len(req.AccountGroupIDs) != 0 {
		filters[pendingLinkedTransferDDL.ColumnAccountGroupID] = req.AccountGroupIDs
	}
	if len(req.TargetAccountIDs) != 0 {
		filters[pendingLinkedTransferDDL.ColumnTargetAccountID] = req.TargetAccountIDs
	}
	if req.Status != nil {
		filters[pendingLinkedTransferDDL.ColumnStatus] = *req.Status
	}

	// Проверяем, что хоть один фильтр был передан
	if len(filters) == 0 {
		return transfers, errors.BadRequest.New("No filters").WithContextParams(ctx)
	}

	if err = r.db.Select(ctx, &transfers, sq.
		Select("*").
		From(pendingLinkedTransferDDL.Table).
		Where(filters),
	); err != nil {
		return transfers, err
	}

	return transfers, nil
}
