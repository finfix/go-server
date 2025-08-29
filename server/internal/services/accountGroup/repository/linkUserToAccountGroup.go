package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *AccountGroupRepository) LinkUserToAccountGroup(ctx context.Context, userID, accountGroupID uuid.UUID) error {
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
