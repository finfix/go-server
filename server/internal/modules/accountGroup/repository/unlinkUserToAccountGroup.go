package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *AccountGroupRepository) UnlinkUserFromAccountGroup(ctx context.Context, userID, accountGroupID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "UnlinkUserFromAccountGroup")
	defer span.End()

	// Исполняем запрос на разрыв связи пользователя с группой счетов
	return r.db.Exec(ctx, sq.
		Delete("coin.users_to_account_groups").
		Where(sq.And{
			sq.Eq{"user_id": userID},
			sq.Eq{"account_group_id": accountGroupID},
		}),
	)
}
