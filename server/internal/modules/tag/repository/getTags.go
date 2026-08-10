package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/modules/accountGroup/repository/accountGroupDDL"
	"server/internal/modules/tag/repository/tagDDL"
	"server/internal/utils/errors"

	"server/internal/modules/tag/model"
)

// GetTags возвращает все подкатегории по фильтрам
func (r *TagRepository) GetTags(ctx context.Context, req model.GetTagsReq) (tags []model.Tag, err error) {
	ctx, span := tracer.Start(ctx, "GetTags")
	defer span.End()

	filtersEq := make(sq.Eq)

	if len(req.AccountGroupIDs) > 0 {
		filtersEq[tagDDL.WithPrefix(tagDDL.ColumnAccountGroupID)] = req.AccountGroupIDs
	}
	if len(req.IDs) > 0 {
		filtersEq[tagDDL.WithPrefix(tagDDL.ColumnID)] = req.IDs
	}

	// Проверяем, что есть фильтры
	if len(filtersEq) == 0 {
		return nil, errors.BadRequest.New("No filters specified").WithContextParams(ctx)
	}

	// Исключаем удаленные подкатегории и подкатегории из удаленных групп счетов
	filtersEq[tagDDL.WithPrefix(tagDDL.ColumnIsDeleted)] = false
	filtersEq[accountGroupDDL.WithPrefix(accountGroupDDL.ColumnIsDeleted)] = false

	// Получаем подкатегории
	return tags, r.db.Select(ctx, &tags, sq.
		Select(tagDDL.WithPrefix(ddlHelper.SelectAll)).
		From(tagDDL.TableWithAlias).
		Join(ddlHelper.BuildJoin(
			accountGroupDDL.TableNameWithAlias,
			accountGroupDDL.WithPrefix(accountGroupDDL.ColumnID),
			tagDDL.WithPrefix(tagDDL.ColumnAccountGroupID),
		)).
		Where(filtersEq),
	)
}
