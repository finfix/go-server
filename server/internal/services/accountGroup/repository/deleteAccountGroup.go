package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/services/accountGroup/repository/accountGroupDDL"
)

// DeleteAccountGroup удаляет группу счетов
func (r *AccountGroupRepository) DeleteAccountGroup(ctx context.Context, id uint32) error {
	ctx, span := tracer.Start(ctx, "DeleteAccountGroup")
	defer span.End()

	// Исполняем запрос на удаление группы счетов
	return r.db.Exec(ctx, sq.
		Delete(accountGroupDDL.TableName).
		Where(sq.Eq{accountGroupDDL.ColumnID: id}),
	)
}
