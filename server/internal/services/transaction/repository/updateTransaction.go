package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/services/transaction/repository/transactionDDL"
	"server/internal/utils/errors"

	"server/internal/services/transaction/model"
)

// UpdateTransaction редактирует транзакцию
func (r *TransactionRepository) UpdateTransaction(ctx context.Context, fields model.UpdateTransactionReq) error {
	ctx, span := tracer.Start(ctx, "UpdateTransaction")
	defer span.End()

	updates := make(map[string]any)

	// Добавляем в запрос поля, которые нужно изменить
	if fields.IsExecuted != nil {
		updates[transactionDDL.ColumnIsExecuted] = *fields.IsExecuted
	}
	if fields.AccountFromID != nil {
		updates[transactionDDL.ColumnAccountFromID] = *fields.AccountFromID
	}
	if fields.AccountToID != nil {
		updates[transactionDDL.ColumnAccountToID] = *fields.AccountToID
	}
	if fields.AmountFrom != nil {
		updates[transactionDDL.ColumnAmountFrom] = *fields.AmountFrom
	}
	if fields.AmountTo != nil {
		updates[transactionDDL.ColumnAmountTo] = *fields.AmountTo
	}
	if fields.DateTransaction != nil {
		updates[transactionDDL.ColumnDate] = *fields.DateTransaction
	}
	if fields.AccountingInCharts != nil {
		updates[transactionDDL.ColumnAccountingInCharts] = *fields.AccountingInCharts
	}
	if fields.Note != nil {
		updates[transactionDDL.ColumnNote] = *fields.Note
	}

	// Проверяем, что есть поля для обновления
	if len(updates) == 0 {
		if fields.TagIDs != nil {
			return nil
		}
		return errors.BadRequest.New("No fields to update").WithContextParams(ctx)
	}

	// Создаем транзакцию
	return r.db.Exec(ctx, sq.
		Update(transactionDDL.Table).
		SetMap(updates).
		Where(sq.Eq{transactionDDL.ColumnID: fields.ID}),
	)
}
