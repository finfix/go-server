package repository

import (
	"context"

	"pkg/ddlHelper"
	"server/internal/services/accountGroup/repository/accountGroupDDL"

	sq "github.com/Masterminds/squirrel"
)

// ChangeSerialNumbers вставляет группу счетов на новое место
func (r *AccountGroupRepository) ChangeSerialNumbers(ctx context.Context, oldValue, newValue uint32) error {
	ctx, span := tracer.Start(ctx, "ChangeSerialNumbers")
	defer span.End()

	// Формируем первичный запрос
	q := sq.Update(accountGroupDDL.TableName)

	// Дополняем запрос в зависимости от того, в какую сторону двигаем элемент
	if newValue < oldValue {
		q = q.
			Set(accountGroupDDL.ColumnSerialNumber, sq.Expr(ddlHelper.Plus(accountGroupDDL.ColumnSerialNumber, 1))).
			Where(sq.And{
				sq.GtOrEq{accountGroupDDL.ColumnSerialNumber: newValue},
				sq.Lt{accountGroupDDL.ColumnSerialNumber: oldValue},
			})
	} else {
		q = q.
			Set(accountGroupDDL.ColumnSerialNumber, sq.Expr(ddlHelper.Minus(accountGroupDDL.ColumnSerialNumber, 1))).
			Where(sq.And{
				sq.Gt{accountGroupDDL.ColumnSerialNumber: oldValue},
				sq.LtOrEq{accountGroupDDL.ColumnSerialNumber: newValue},
			})
	}

	// Выполняем запрос
	return r.db.Exec(ctx, q)
}
