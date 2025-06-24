package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/services/tag/repository/tagToTransactionDDL"
)

// UnlinkTagsFromTransaction отвязывает подкатегории от транзакции
func (r *TagRepository) UnlinkTagsFromTransaction(ctx context.Context, tagIDs []uint32, transactionID uint32) error {
	ctx, span := tracer.Start(ctx, "UnlinkTagsFromTransaction")
	defer span.End()

	filtersEq := make(sq.Eq)

	filtersEq[tagToTransactionDDL.ColumnTagID] = tagIDs
	filtersEq[tagToTransactionDDL.ColumnTransactionID] = transactionID

	// Удаляем связь между подкатегориями и транзакцией
	return r.db.Exec(ctx, sq.
		Delete(tagToTransactionDDL.Table).
		Where(filtersEq),
	)
}
