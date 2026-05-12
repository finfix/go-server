package repository

import (
	"context"

	userModel "server/internal/modules/user/model"
	"server/internal/modules/user/repository/userDDL"

	sq "github.com/Masterminds/squirrel"
)

// CreateUser Создает нового пользователя
func (r *UserRepository) CreateUser(ctx context.Context, user userModel.CreateReq) error {
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()

	return r.db.Exec(ctx, sq.
		Insert(`coin.users`).
		SetMap(map[string]any{
			userDDL.ColumnID:              user.ID,
			userDDL.ColumnName:            user.Name,
			userDDL.ColumnEmail:           user.Email,
			userDDL.ColumnPasswordHash:    user.PasswordHash,
			userDDL.ColumnTimeCreate:      user.TimeCreate,
			userDDL.ColumnDefaultCurrency: user.DefaultCurrency,
			userDDL.ColumnPasswordSalt:    user.PasswordSalt,
			userDDL.ColumnIsAdmin:         user.IsAdmin,
		}),
	)
}
