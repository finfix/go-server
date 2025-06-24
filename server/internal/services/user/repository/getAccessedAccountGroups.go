package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"pkg/log"
	"server/internal/services/accountGroup/repository/accountGroupDDL"
	"server/internal/services/user/model"
	"server/internal/services/user/repository/userDDL"
	"server/internal/services/user/repository/userToAccountGroupDDL"
)

// GetAccessedAccountGroups возвращает доступы пользователей к группам счетов
func (r *UserRepository) GetAccessedAccountGroups(ctx context.Context, userID uint32) (accessedAccountGroupIDs []uint32, err error) {
	ctx, span := tracer.Start(ctx, "GetAccessedAccountGroups")
	defer span.End()

	// Проверяем наличие данных в кэше
	var ok bool
	if accessedAccountGroupIDs, ok = r.accessedAccountGroupIDsCache.Get(userID); ok {
		return accessedAccountGroupIDs, nil
	}

	log.WithContextParams(ctx).Info("Обновляем кэш с доступами пользователей к группам счетов")

	var usersToAccountGroups []model.UserToAccountGroup

	// Получаем связки между пользователями и группами счетов
	if err = r.db.Select(ctx, &usersToAccountGroups, sq.
		Select(
			ddlHelper.As(
				userDDL.WithPrefix(userDDL.ColumnID),
				"user_id",
			),
			ddlHelper.As(
				accountGroupDDL.WithPrefix(accountGroupDDL.ColumnID),
				"account_group_id",
			),
		).
		From(accountGroupDDL.TableNameWithAlias).
		Join(ddlHelper.BuildJoin(
			userToAccountGroupDDL.TableWithAlias,
			userToAccountGroupDDL.WithPrefix(userToAccountGroupDDL.ColumnAccountGroupID),
			accountGroupDDL.WithPrefix(accountGroupDDL.ColumnID),
		)).
		Join(ddlHelper.BuildJoin(
			userDDL.TableWithAlias,
			userDDL.WithPrefix(userDDL.ColumnID),
			userToAccountGroupDDL.WithPrefix(userToAccountGroupDDL.ColumnUserID),
		)),
	); err != nil {
		return nil, err
	}

	// Формируем мапу userID - []accountGroupID
	usersToAccountsGroups := make(map[uint32][]uint32)

	// Проходимся по каждой связке пользователя и группы счетов
	for _, userToAccountGroup := range usersToAccountGroups {
		accessedAccountGroups := usersToAccountsGroups[userToAccountGroup.UserID]
		accessedAccountGroups = append(accessedAccountGroups, userToAccountGroup.AccountGroupID)
		usersToAccountsGroups[userToAccountGroup.UserID] = accessedAccountGroups
	}

	// Проходимся по каждому пользователю и его группам счетов
	for userID, accountGroups := range usersToAccountsGroups {

		// Добавляем запись в кэш
		r.accessedAccountGroupIDsCache.Set(userID, accountGroups)
	}

	return usersToAccountsGroups[userID], nil
}
