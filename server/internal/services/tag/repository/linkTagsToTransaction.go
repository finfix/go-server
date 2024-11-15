package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/services/tag/repository/tagToTransactionDDL"
)

// LinkTagsToTransaction привязывает подкатегории к транзакции
func (r *TagRepository) LinkTagsToTransaction(ctx context.Context, tagIDs []uint32, transactionID uint32) error {
	ctx, span := tracer.Start(ctx, "LinkTagsToTransaction")
	defer span.End()

	q := sq.
		Insert(tagToTransactionDDL.Table).
		Columns(tagToTransactionDDL.ColumnTransactionID, tagToTransactionDDL.ColumnTagID)

	for _, tagID := range tagIDs {
		q = q.Values(transactionID, tagID)
	}

	return r.db.Exec(ctx, q)
}
