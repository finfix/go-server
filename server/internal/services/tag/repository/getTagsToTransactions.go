package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/services/accountGroup/repository/accountGroupDDL"
	"server/internal/services/tag/model"
	"server/internal/services/tag/repository/tagDDL"
	"server/internal/services/tag/repository/tagToTransactionDDL"
	"server/internal/services/transaction/repository/transactionDDL"
)

// GetTagsToTransactions возвращает все связи между подкатегориями и транзакциями
func (r *TagRepository) GetTagsToTransactions(ctx context.Context, req model.GetTagsToTransactionsReq) (res []model.TagToTransaction, err error) {
	ctx, span := tracer.Start(ctx, "GetTagsToTransactions")
	defer span.End()

	// Формируем первичный запрос
	q := sq.
		Select("*").
		From("coin.tags_to_transaction ttt")

	// Фильтрация по переданным группам счетов
	if len(req.AccountGroupIDs) != 0 {
		q = q.
			Join(ddlHelper.BuildJoin(
				tagDDL.TableWithAlias,
				tagDDL.WithPrefix(tagDDL.ColumnID),
				tagToTransactionDDL.WithPrefix(tagToTransactionDDL.ColumnTagID),
			)).
			Join(ddlHelper.BuildJoin(
				accountGroupDDL.TableNameWithAlias,
				accountGroupDDL.WithPrefix(accountGroupDDL.ColumnID),
				tagDDL.WithPrefix(tagDDL.ColumnAccountGroupID),
			)).
			Where(sq.Eq{tagDDL.WithPrefix(tagDDL.ColumnAccountGroupID): req.AccountGroupIDs})
	}

	// Фильтрация по переданным транзакциям
	if len(req.TransactionIDs) != 0 {
		q = q.
			Join(ddlHelper.BuildJoin(
				transactionDDL.TableWithAlias,
				transactionDDL.WithPrefix(transactionDDL.ColumnID),
				tagToTransactionDDL.WithPrefix(tagToTransactionDDL.ColumnTransactionID),
			)).
			Where(sq.Eq{tagToTransactionDDL.WithPrefix(tagToTransactionDDL.ColumnTransactionID): req.TransactionIDs})
	}

	// Выполняем запрос
	return res, r.db.Select(ctx, &res, q)
}
