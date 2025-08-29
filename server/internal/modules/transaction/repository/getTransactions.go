package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/modules/account/repository/accountDDL"
	"server/internal/modules/transaction/model"
	"server/internal/modules/transaction/repository/transactionDDL"
)

// GetTransactions возвращает все транзакции по фильтрам
func (r *TransactionRepository) GetTransactions(ctx context.Context, req model.GetTransactionsReq) (transactions []model.Transaction, err error) {
	ctx, span := tracer.Start(ctx, "GetTransactions")
	defer span.End()

	accountsFromPrefix, accountsToPrefix := "a1", "a2"

	q := sq.
		Select(transactionDDL.WithPrefix(ddlHelper.SelectAll)).
		From(transactionDDL.TableWithAlias)

	// Добавляем фильтры
	if len(req.AccountGroupIDs) != 0 {
		q = q.
			Join(ddlHelper.BuildJoin(
				ddlHelper.WithCustomAlias(accountDDL.Table, accountsFromPrefix),
				ddlHelper.WithCustomPrefix(accountDDL.ColumnID, accountsFromPrefix),
				transactionDDL.WithPrefix(transactionDDL.ColumnAccountFromID),
			)).
			Join(ddlHelper.BuildJoin(
				ddlHelper.WithCustomAlias(accountDDL.Table, accountsToPrefix),
				ddlHelper.WithCustomPrefix(accountDDL.ColumnID, accountsToPrefix),
				transactionDDL.WithPrefix(transactionDDL.ColumnAccountToID),
			)).
			Where(sq.Eq{
				ddlHelper.WithCustomPrefix(accountDDL.ColumnAccountGroupID, accountsFromPrefix): req.AccountGroupIDs,
				ddlHelper.WithCustomPrefix(accountDDL.ColumnAccountGroupID, accountsToPrefix):   req.AccountGroupIDs,
			})
	}
	if len(req.IDs) != 0 {
		q = q.Where(sq.Eq{transactionDDL.ColumnID: req.IDs})
	}
	if req.AccountID != nil {
		q = q.Where(sq.Or{
			sq.Eq{transactionDDL.WithPrefix(transactionDDL.ColumnAccountFromID): *req.AccountID},
			sq.Eq{transactionDDL.WithPrefix(transactionDDL.ColumnAccountToID): *req.AccountID},
		})
	}
	if req.Type != nil {
		q = q.Where(sq.Eq{transactionDDL.WithPrefix(transactionDDL.ColumnType): *req.Type})
	}
	if req.DateFrom != nil {
		q = q.Where(sq.GtOrEq{transactionDDL.WithPrefix(transactionDDL.ColumnDate): req.DateFrom})
	}
	if req.DateTo != nil {
		q = q.Where(sq.LtOrEq{transactionDDL.WithPrefix(transactionDDL.ColumnDate): req.DateTo})
	}

	q = q.OrderBy(
		ddlHelper.Desc(transactionDDL.WithPrefix(transactionDDL.ColumnDate)),
		ddlHelper.Desc(transactionDDL.WithPrefix(transactionDDL.ColumnDatetimeCreate)),
	)

	if req.Limit != nil {
		q = q.Limit(uint64(*req.Limit))
	}
	if req.Offset != nil {
		q = q.Offset(uint64(*req.Offset))
	}

	// Получаем транзакции
	if err = r.db.Select(ctx, &transactions, q); err != nil {
		return nil, err
	}

	return transactions, nil
}
