package service

import (
	"context"

	"server/internal/services/tag/model"
	"server/internal/utils/errors"
)

func (s *TagService) GetTags(ctx context.Context, filters model.GetTagsReq) (tags []model.Tag, err error) {
	ctx, span := tracer.Start(ctx, "GetTags")
	defer span.End()

	// Проверяем доступ пользователя к группам счетов
	if filters.AccountGroupIDs != nil {
		if err = s.accountGroupService.CheckAccess(ctx, filters.Necessary.UserID, filters.AccountGroupIDs); err != nil {
			return nil, err
		}
	} else {
		filters.AccountGroupIDs, err = s.userService.GetAccessedAccountGroups(ctx, filters.Necessary.UserID)
		if err != nil {
			return nil, err
		}
		if len(filters.AccountGroupIDs) == 0 {
			return nil, errors.Forbidden.New("У пользователя нет доступа к группам счетов").WithContextParams(ctx)
		}
	}

	// Получаем все подкатегории
	if tags, err = s.tagRepository.GetTags(ctx, filters); err != nil {
		return nil, err
	}

	// Заполняем массив ID транзакций
	tagIDs := make([]uint32, len(tags))
	for i, tag := range tags {
		tagIDs[i] = tag.ID
	}

	return tags, nil
}
