package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/services/transaction/repository/transactionDDL"
	"server/internal/utils/errors"
)

// CheckAccess проверяет, имеет ли набор групп счетов пользователя доступ к указанным идентификаторам транзакций
func (r *TransactionRepository) CheckAccess(ctx context.Context, accountGroupIDs, transactionIDs []uint32) error {
	ctx, span := tracer.Start(ctx, "CheckAccess")
	defer span.End()

	// Получаем все доступные транзакции по группам счетов и перечисленным транзакциям
	rows, err := r.db.Query(ctx, sq.
		Select(transactionDDL.ColumnID).
		From(transactionDDL.TableWithAlias).
		Where(sq.Eq{
			transactionDDL.ColumnAccountGroupID: accountGroupIDs,
			transactionDDL.ColumnID:             transactionIDs,
		}))
	if err != nil {
		return err
	}

	// Формируем мапу доступных транзакций
	accessedTransactionIDs := make(map[uint32]struct{})

	// Проходимся по каждой доступной транзакции
	for rows.Next() {

		// Считываем ID транзакции
		var transactionID uint32
		if err = rows.Scan(&transactionID); err != nil {
			return err
		}

		// Добавляем ID транзакции в мапу
		accessedTransactionIDs[transactionID] = struct{}{}
	}

	if len(accessedTransactionIDs) == 0 {
		return errors.Forbidden.New("You don't have access to any of the requested transactions").
			WithContextParams(ctx).
			WithParams("AccountGroupIDs", accountGroupIDs, "TransactionIDs", transactionIDs)
	}

	// Проходимся по каждой запрашиваемой транзакции
	for _, transactionID := range transactionIDs {

		// Если счета нет в мапе доступных транзакций, возвращаем ошибку
		if _, ok := accessedTransactionIDs[transactionID]; !ok {
			return errors.Forbidden.New("You don't have access to transaction").
				WithContextParams(ctx).
				WithParams(
					"AccountGroupIDs", accountGroupIDs,
					"TransactionID", transactionID,
				).
				SkipPreviousCaller()
		}
	}

	return nil
}
