package service

import (
	"context"

	"server/internal/services/tag/model"
)

// CreateTag создает новую подкатегорию
func (s *TagService) CreateTag(ctx context.Context, tag model.CreateTagReq) (id uint32, err error) {
	ctx, span := tracer.Start(ctx, "CreateTag")
	defer span.End()

	// Проверяем доступ пользователя к группам счетов
	if err = s.accountGroupService.CheckAccess(ctx, tag.Necessary.UserID, []uint32{tag.AccountGroupID}); err != nil {
		return id, err
	}

	// Создаем подкатегорию
	return s.tagRepository.CreateTag(ctx, tag.ConvertToRepoReq())
}
