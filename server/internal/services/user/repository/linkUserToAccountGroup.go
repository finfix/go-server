package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/services/user/repository/userToAccountGroupDDL"
)

func (r *UserRepository) LinkUserToAccountGroup(ctx context.Context, userID uint32, accountGroupID uint32) error {
	ctx, span := tracer.Start(ctx, "LinkUserToAccountGroup")
	defer span.End()

	return r.db.Exec(ctx, sq.
		Insert(userToAccountGroupDDL.Table).
		SetMap(map[string]any{
			userToAccountGroupDDL.ColumnUserID:         userID,
			userToAccountGroupDDL.ColumnAccountGroupID: accountGroupID,
		}),
	)
}
