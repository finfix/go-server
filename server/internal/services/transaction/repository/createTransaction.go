package repository

import (
	"context"

	"github.com/google/uuid"
	sq "github.com/Masterminds/squirrel"

	"server/internal/services/transaction/repository/model"
	"server/internal/services/transaction/repository/transactionDDL"
)

// CreateTransaction создает новую транзакцию
func (r *TransactionRepository) CreateTransaction(ctx context.Context, req model.CreateTransactionReq) (id uuid.UUID, err error) {
	ctx, span := tracer.Start(ctx, "CreateTransaction")
	defer span.End()

	// Создаем транзакцию
	return r.db.ExecWithLastUUID(ctx, sq.Insert(`coin.transactions`).
		SetMap(map[string]any{
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
}
