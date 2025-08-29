package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/accountGroup/repository/accountGroupDDL"
)

// DeleteAccountGroup удаляет группу счетов
func (r *AccountGroupRepository) DeleteAccountGroup(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "DeleteAccountGroup")
	defer span.End()

	// Исполняем запрос на удаление группы счетов
	return r.db.Exec(ctx, sq.
		Delete(accountGroupDDL.TableName).
		Where(sq.Eq{accountGroupDDL.ColumnID: id}),
	)
}
