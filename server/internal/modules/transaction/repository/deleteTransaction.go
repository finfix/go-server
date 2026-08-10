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

	// Помечаем транзакцию как удаленную
	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(transactionDDL.Table).
		Set(transactionDDL.ColumnIsDeleted, true).
		Where(sq.Eq{transactionDDL.ColumnID: id, transactionDDL.ColumnIsDeleted: false}),
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
