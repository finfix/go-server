package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"pkg/errors"
	"server/internal/services/tag/repository/tagDDL"

	"server/internal/services/tag/model"
)

// GetTags возвращает все подкатегории по фильтрам
func (r *TagRepository) GetTags(ctx context.Context, req model.GetTagsReq) (tags []model.Tag, err error) {
	ctx, span := tracer.Start(ctx, "GetTags")
	defer span.End()

	filtersEq := make(sq.Eq)

	if len(req.AccountGroupIDs) > 0 {
		filtersEq[tagDDL.ColumnAccountGroupID] = req.AccountGroupIDs
	}

	// Проверяем, что есть фильтры
	if len(filtersEq) == 0 {
		return nil, errors.BadRequest.New(ctx, "No filters specified")
	}

	// Получаем подкатегории
	return tags, r.db.Select(ctx, &tags, sq.
		Select(ddlHelper.SelectAll).
		From(tagDDL.Table).
		Where(filtersEq),
	)
}
