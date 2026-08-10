package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/tag/repository/tagDDL"
	"server/internal/utils/errors"
)

// DeleteTag удаляет подкатегорию
func (r *TagRepository) DeleteTag(ctx context.Context, id, userID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "DeleteTag")
	defer span.End()

	// Помечаем подкатегорию как удаленную
	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(tagDDL.Table).
		Set(tagDDL.ColumnIsDeleted, true).
		Where(sq.Eq{tagDDL.ColumnID: id, tagDDL.ColumnIsDeleted: false}),
	)
	if err != nil {
		return err
	}

	// Проверяем, что в базе данных что-то изменилось
	if rows == 0 {
		return errors.NotFound.New("No such model found for model").
			WithContextParams(ctx).
			WithParams(
				"UserID", userID,
				"TagID", id,
			)
	}

	return nil
}
