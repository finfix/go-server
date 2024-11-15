package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	userModel "server/internal/services/user/model"
	"server/internal/services/user/repository/userDDL"
)

// CreateUser Создает нового пользователя
func (r *UserRepository) CreateUser(ctx context.Context, user userModel.CreateReq) (uint32, error) {
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()

	return r.db.ExecWithLastInsertID(ctx, sq.
		Insert(`coin.users`).
		SetMap(map[string]any{
			userDDL.ColumnName:            user.Name,
			userDDL.ColumnEmail:           user.Email,
			userDDL.ColumnPasswordHash:    user.PasswordHash,
			userDDL.ColumnTimeCreate:      user.TimeCreate,
			userDDL.ColumnDefaultCurrency: user.DefaultCurrency,
			userDDL.ColumnPasswordSalt:    user.PasswordSalt,
			userDDL.ColumnIsAdmin:         false,
		}),
	)
}
