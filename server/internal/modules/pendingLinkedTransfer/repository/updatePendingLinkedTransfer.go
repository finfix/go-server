package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/pendingLinkedTransfer/repository/model"
	"server/internal/modules/pendingLinkedTransfer/repository/pendingLinkedTransferDDL"
	"server/internal/utils/errors"
)

// UpdatePendingLinkedTransfer обновляет статус переноса (patch — только переданные поля)
func (r *PendingLinkedTransferRepository) UpdatePendingLinkedTransfer(ctx context.Context, id uuid.UUID, req model.UpdatePendingLinkedTransferReq) error {
	ctx, span := tracer.Start(ctx, "updatePendingLinkedTransfer")
	defer span.End()

	updates := make(map[string]any)
	if req.Status != nil {
		updates[pendingLinkedTransferDDL.ColumnStatus] = *req.Status
	}

	if len(updates) == 0 {
		return errors.BadRequest.New("No fields to update").WithContextParams(ctx)
	}

	return r.db.Exec(ctx, sq.
		Update(pendingLinkedTransferDDL.Table).
		SetMap(updates).
		Where(sq.Eq{pendingLinkedTransferDDL.ColumnID: id}),
	)
}
