package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/services/tag/repository/tagToTransactionDDL"
)

// UnlinkTagsFromTransaction отвязывает подкатегории от транзакции
func (r *TagRepository) UnlinkTagsFromTransaction(ctx context.Context, tagIDs []uuid.UUID, transactionID uuid.UUID) error {
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
