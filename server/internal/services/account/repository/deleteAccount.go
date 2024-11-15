package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/services/account/repository/accountDDL"
)

// DeleteAccount удаляет счет
func (r *AccountRepository) DeleteAccount(ctx context.Context, id uint32) error {
	ctx, span := tracer.Start(ctx, "deleteAccount")
	defer span.End()

	// Исполняем запрос на удаление счета
	return r.db.Exec(ctx, sq.
		Delete(accountDDL.Table).
		Where(sq.Eq{accountDDL.ColumnID: id}),
	)
}
