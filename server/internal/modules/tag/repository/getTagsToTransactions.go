package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/modules/accountGroup/repository/accountGroupDDL"
	"server/internal/modules/tag/model"
	"server/internal/modules/tag/repository/tagDDL"
	"server/internal/modules/tag/repository/tagToTransactionDDL"
	"server/internal/modules/transaction/repository/transactionDDL"
)

// GetTagsToTransactions возвращает все связи между подкатегориями и транзакциями
func (r *TagRepository) GetTagsToTransactions(ctx context.Context, req model.GetTagsToTransactionsReq) (res []model.TagToTransaction, err error) {
	ctx, span := tracer.Start(ctx, "GetTagsToTransactions")
	defer span.End()

	// Формируем первичный запрос, исключая связи с удаленными подкатегориями, транзакциями и группами счетов
	q := sq.
		Select(tagToTransactionDDL.WithPrefix(ddlHelper.SelectAll)).
		From(tagToTransactionDDL.TableWithAlias).
		Join(ddlHelper.BuildJoin(
			tagDDL.TableWithAlias,
			tagDDL.WithPrefix(tagDDL.ColumnID),
			tagToTransactionDDL.WithPrefix(tagToTransactionDDL.ColumnTagID),
		)).
		Join(ddlHelper.BuildJoin(
			transactionDDL.TableWithAlias,
			transactionDDL.WithPrefix(transactionDDL.ColumnID),
			tagToTransactionDDL.WithPrefix(tagToTransactionDDL.ColumnTransactionID),
		)).
		Join(ddlHelper.BuildJoin(
			accountGroupDDL.TableNameWithAlias,
			accountGroupDDL.WithPrefix(accountGroupDDL.ColumnID),
			tagDDL.WithPrefix(tagDDL.ColumnAccountGroupID),
		)).
		Where(sq.Eq{
			tagDDL.WithPrefix(tagDDL.ColumnIsDeleted):                   false,
			transactionDDL.WithPrefix(transactionDDL.ColumnIsDeleted):   false,
			accountGroupDDL.WithPrefix(accountGroupDDL.ColumnIsDeleted): false,
		})

	// Фильтрация по переданным группам счетов
	if len(req.AccountGroupIDs) != 0 {
		q = q.Where(sq.Eq{tagDDL.WithPrefix(tagDDL.ColumnAccountGroupID): req.AccountGroupIDs})
	}

	// Фильтрация по переданным транзакциям
	if len(req.TransactionIDs) != 0 {
		q = q.Where(sq.Eq{tagToTransactionDDL.WithPrefix(tagToTransactionDDL.ColumnTransactionID): req.TransactionIDs})
	}

	// Выполняем запрос
	return res, r.db.Select(ctx, &res, q)
}
