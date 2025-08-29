package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	tagRepoModel "server/internal/modules/tag/repository/model"
	"server/internal/modules/tag/repository/tagDDL"
)

// CreateTag создает новую подкатегорию
func (r *TagRepository) CreateTag(ctx context.Context, req tagRepoModel.CreateTagReq) (id uuid.UUID, err error) {
	ctx, span := tracer.Start(ctx, "CreateTag")
	defer span.End()

	// Создаем подкатегорию
	return r.db.ExecWithLastUUID(ctx, sq.
		Insert(tagDDL.Table).
		SetMap(map[string]any{
			tagDDL.ColumnName:            req.Name,
			tagDDL.ColumnAccountGroupID:  req.AccountGroupID,
			tagDDL.ColumnCreatedByUserID: req.CreatedByUserID,
			tagDDL.ColumnDatetimeCreate:  req.DatetimeCreate,
		}),
	)
}
