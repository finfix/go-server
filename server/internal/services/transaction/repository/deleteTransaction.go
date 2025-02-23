package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/errors"
	"server/internal/services/transaction/repository/transactionDDL"
)

// DeleteTransaction удаляет транзакцию
func (r *TransactionRepository) DeleteTransaction(ctx context.Context, id, userID uint32) error {
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
		return errors.NotFound.New(ctx, "No such model found for model",
			errors.ParamsOption("UserID", userID, "TransactionID", id),
		)
	}

	return nil
}
