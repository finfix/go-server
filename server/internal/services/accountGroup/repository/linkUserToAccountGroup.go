package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *AccountGroupRepository) LinkUserToAccountGroup(ctx context.Context, userID, accountGroupID uint32) error {
	ctx, span := tracer.Start(ctx, "LinkUserToAccountGroup")
	defer span.End()

	// Исполняем запрос на связывание пользователя с группой счетов
	return r.db.Exec(ctx, sq.
		Insert("coin.users_to_account_groups").
		SetMap(map[string]any{
			"user_id":          userID,
			"account_group_id": accountGroupID,
		}),
	)
}
