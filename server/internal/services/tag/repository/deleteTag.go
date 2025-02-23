package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/errors"
	"server/internal/services/tag/repository/tagDDL"
)

// DeleteTag удаляет подкатегорию
func (r *TagRepository) DeleteTag(ctx context.Context, id, userID uint32) error {
	ctx, span := tracer.Start(ctx, "DeleteTag")
	defer span.End()

	// Удаляем подкатегорию
	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Delete(tagDDL.Table).
		Where(sq.Eq{tagDDL.ColumnID: id}),
	)
	if err != nil {
		return err
	}

	// Проверяем, что в базе данных что-то изменилось
	if rows == 0 {
		return errors.NotFound.New(ctx, "No such model found for model",
			errors.ParamsOption("UserID", userID, "TagID", id),
		)
	}

	return nil
}
