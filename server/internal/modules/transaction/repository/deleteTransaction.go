package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/transaction/repository/transactionDDL"
	"server/internal/utils/errors"
)

// DeleteTransaction удаляет транзакцию
func (r *TransactionRepository) DeleteTransaction(ctx context.Context, id, userID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "DeleteTransaction")
	defer span.End()

	// Удаляем транзакцию
	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Delete(transactionDDL.Table).
		Where(sq.Eq{transactionDDL.ColumnID: id}),
	)
	if err != nil {
		return err
	}

	// Проверяем, что в базе данных что-то изменилось
	if rows == 0 {
		return errors.NotFound.New("No such model found for model").
			WithContextParams(ctx).
			WithParams("UserID", userID, "TransactionID", id)
	}

	return nil
}
