package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/transaction/repository/model"
	"server/internal/modules/transaction/repository/transactionDDL"
)

// CreateTransaction создает новую транзакцию
func (r *TransactionRepository) CreateTransaction(ctx context.Context, req model.CreateTransactionReq) (id uuid.UUID, err error) {
	ctx, span := tracer.Start(ctx, "CreateTransaction")
	defer span.End()

	// Используем ID из запроса
	id = req.ID

	// Создаем транзакцию
	err = r.db.Exec(ctx, sq.Insert(`coin.transactions`).
		SetMap(map[string]any{
			transactionDDL.ColumnID:                 req.ID,
			transactionDDL.ColumnType:               req.Type,
			transactionDDL.ColumnDate:               req.DateTransaction,
			transactionDDL.ColumnAccountFromID:      req.AccountFromID,
			transactionDDL.ColumnAccountToID:        req.AccountToID,
			transactionDDL.ColumnAmountFrom:         req.AmountFrom,
			transactionDDL.ColumnAmountTo:           req.AmountTo,
			transactionDDL.ColumnNote:               req.Note,
			transactionDDL.ColumnIsExecuted:         req.IsExecuted,
			transactionDDL.ColumnDatetimeCreate:     req.DatetimeCreate,
			transactionDDL.ColumnCreatedByUserID:    req.CreatedByUserID,
			transactionDDL.ColumnAccountingInCharts: req.AccountingInCharts,
			transactionDDL.ColumnAccountGroupID:     req.AccountGroupID,
		}),
	)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
