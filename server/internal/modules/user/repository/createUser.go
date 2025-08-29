package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	userModel "server/internal/modules/user/model"
	"server/internal/modules/user/repository/userDDL"
)

// CreateUser Создает нового пользователя
func (r *UserRepository) CreateUser(ctx context.Context, user userModel.CreateReq) (uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()

	return r.db.ExecWithLastUUID(ctx, sq.
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
