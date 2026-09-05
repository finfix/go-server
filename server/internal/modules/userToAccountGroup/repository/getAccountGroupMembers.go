package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/userToAccountGroup/repository/userToAccountGroupDDL"
)

// GetAccountGroupMembers возвращает всех пользователей, имеющих доступ к группе счетов —
// обратная выборка к GetAccessedAccountGroups (там пользователь → группы, здесь группа →
// пользователи). Используется, чтобы понять, кого будить через SyncNotifier при изменении
// сущности в этой группе.
func (r *UserToAccountGroupRepository) GetAccountGroupMembers(ctx context.Context, accountGroupID uuid.UUID) (userIDs []uuid.UUID, err error) {
	ctx, span := tracer.Start(ctx, "GetAccountGroupMembers")
	defer span.End()

	if err = r.db.Select(ctx, &userIDs, sq.
		Select(userToAccountGroupDDL.ColumnUserID).
		From(userToAccountGroupDDL.Table).
		Where(sq.Eq{userToAccountGroupDDL.ColumnAccountGroupID: accountGroupID}),
	); err != nil {
		return nil, err
	}

	return userIDs, nil
}
