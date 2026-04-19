package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"pkg/ddlHelper"
	"server/internal/modules/accountGroup/repository/accountGroupDDL"
	"server/internal/modules/user/repository/userToAccountGroupDDL"

	"server/internal/modules/accountGroup/model"
)

func (r *AccountGroupRepository) GetAccountGroups(ctx context.Context, req model.GetAccountGroupsReq) (accountGroups []model.AccountGroup, err error) {
	ctx, span := tracer.Start(ctx, "GetAccountGroups")
	defer span.End()

	// Формируем первичный запрос
	q := sq.
		Select(accountGroupDDL.WithPrefix(ddlHelper.SelectAll)).
		From(accountGroupDDL.TableNameWithAlias)

	// Фильтр по группам счетов
	if len(req.AccountGroupIDs) != 0 {
		q = q.Where(sq.Eq{accountGroupDDL.WithPrefix(accountGroupDDL.ColumnID): req.AccountGroupIDs})
	}

	// Фильтр по пользователю
	if req.Necessary.UserID != uuid.Nil {
		q = q.
			Join(ddlHelper.BuildJoin(
				userToAccountGroupDDL.TableWithAlias,
				userToAccountGroupDDL.WithPrefix(userToAccountGroupDDL.ColumnAccountGroupID),
				accountGroupDDL.WithPrefix(accountGroupDDL.ColumnID),
			)).
			Where(sq.Eq{userToAccountGroupDDL.WithPrefix(userToAccountGroupDDL.ColumnUserID): req.Necessary.UserID})
	}

	// Выполняем запрос
	return accountGroups, r.db.Select(ctx, &accountGroups, q)
}
