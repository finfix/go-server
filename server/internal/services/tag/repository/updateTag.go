package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/errors"
	"server/internal/services/tag/repository/tagDDL"

	"server/internal/services/tag/model"
)

// UpdateTag редактирует подкатегорию
func (r *TagRepository) UpdateTag(ctx context.Context, fields model.UpdateTagReq) error {
	ctx, span := tracer.Start(ctx, "UpdateTag")
	defer span.End()

	updates := make(map[string]any)

	// Добавляем в запрос поля, которые нужно изменить
	if fields.Name != nil {
		updates[tagDDL.ColumnName] = *fields.Name
	}

	// Проверяем, что есть поля для обновления
	if len(updates) == 0 {
		return errors.BadRequest.New("No fields to update")
	}

	// Редактируем подкатегорию
	return r.db.Exec(ctx, sq.
		Update(tagDDL.Table).
		SetMap(updates).
		Where(sq.Eq{"id": fields.ID}),
	)
}
