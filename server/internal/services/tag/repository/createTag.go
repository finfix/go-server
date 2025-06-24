package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	tagRepoModel "server/internal/services/tag/repository/model"
	"server/internal/services/tag/repository/tagDDL"
)

// CreateTag создает новую подкатегорию
func (r *TagRepository) CreateTag(ctx context.Context, req tagRepoModel.CreateTagReq) (id uint32, err error) {
	ctx, span := tracer.Start(ctx, "CreateTag")
	defer span.End()

	// Создаем подкатегорию
	return r.db.ExecWithLastInsertID(ctx, sq.
		Insert(tagDDL.Table).
		SetMap(map[string]any{
			tagDDL.ColumnName:            req.Name,
			tagDDL.ColumnAccountGroupID:  req.AccountGroupID,
			tagDDL.ColumnCreatedByUserID: req.CreatedByUserID,
			tagDDL.ColumnDatetimeCreate:  req.DatetimeCreate,
		}),
	)
}
