package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/user/repository/userToAccountGroupDDL"
)

func (r *UserRepository) LinkUserToAccountGroup(ctx context.Context, userID uuid.UUID, accountGroupID uuid.UUID) error {
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
