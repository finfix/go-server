package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/errors"
	userRepoModel "server/internal/services/user/repository/model"
	"server/internal/services/user/repository/userDDL"
)

// UpdateUser редактирует пользователя
func (r *UserRepository) UpdateUser(ctx context.Context, fields userRepoModel.UpdateUserReq) error {
	ctx, span := tracer.Start(ctx, "UpdateUser")
	defer span.End()

	// Изменяем поля пользователя
	updates := make(map[string]any)

	// Добавляем в запрос поля, которые нужно изменить
	if fields.Email != nil {
		updates[userDDL.ColumnEmail] = *fields.Email
	}
	if fields.Name != nil {
		updates[userDDL.ColumnName] = *fields.Name
	}
	if fields.DefaultCurrency != nil {
		updates[userDDL.ColumnDefaultCurrency] = *fields.DefaultCurrency
	}
	if fields.PasswordHash != nil && fields.PasswordSalt != nil {
		updates[userDDL.ColumnPasswordHash] = *fields.PasswordHash
		updates[userDDL.ColumnPasswordSalt] = *fields.PasswordSalt
	}

	if len(updates) == 0 {
		return errors.BadRequest.New(ctx, "No fields to update")
	}

	// Обновляем пользователя
	return r.db.Exec(ctx, sq.
		Update(userDDL.Table).
		SetMap(updates).
		Where(sq.Eq{userDDL.ColumnID: fields.ID}),
	)
}
