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
		return errors.NotFound.New("No such model found for model").
			WithContextParams(ctx).
			WithParams(
				"UserID", userID,
				"TagID", id,
			)
	}

	return nil
}
