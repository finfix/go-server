package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/pendingLinkedTransfer/repository/pendingLinkedTransferDDL"
	"server/internal/utils/errors"
)

// DeletePendingLinkedTransfer удаляет перенос. Жёсткое удаление (без is_deleted) — в отличие от
// транзакций/тегов эта сущность не хранит историю сама по себе, а факт удаления и так попадает
// в аудит-лог через TrackMutation на уровне сервиса.
func (r *PendingLinkedTransferRepository) DeletePendingLinkedTransfer(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "deletePendingLinkedTransfer")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Delete(pendingLinkedTransferDDL.Table).
		Where(sq.Eq{pendingLinkedTransferDDL.ColumnID: id}),
	)
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.NotFound.New("No such model found for model").
			WithContextParams(ctx).
			WithParams("PendingLinkedTransferID", id)
	}

	return nil
}
