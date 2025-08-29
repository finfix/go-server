package repository

import (
	"context"

	tagRepoModel "server/internal/modules/tag/repository/model"
	"server/internal/modules/tag/repository/tagDDL"

	sq "github.com/Masterminds/squirrel"
)

// CreateTag создает новую подкатегорию
func (r *TagRepository) CreateTag(ctx context.Context, req tagRepoModel.CreateTagReq) error {
	ctx, span := tracer.Start(ctx, "CreateTag")
	defer span.End()

	// Создаем подкатегорию
	return r.db.Exec(ctx, sq.
		Insert(tagDDL.Table).
		SetMap(map[string]any{
			tagDDL.ColumnID:              req.ID,
			tagDDL.ColumnName:            req.Name,
			tagDDL.ColumnAccountGroupID:  req.AccountGroupID,
			tagDDL.ColumnCreatedByUserID: req.CreatedByUserID,
			tagDDL.ColumnDatetimeCreate:  req.DatetimeCreate,
		}),
	)
}
