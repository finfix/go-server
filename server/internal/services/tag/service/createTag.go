package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/services/tag/model"
)

// CreateTag создает новую подкатегорию
func (s *TagService) CreateTag(ctx context.Context, tag model.CreateTagReq) (id uuid.UUID, err error) {
	ctx, span := tracer.Start(ctx, "CreateTag")
	defer span.End()

	// Проверяем доступ пользователя к группам счетов
	if err = s.accountGroupService.CheckAccess(ctx, tag.Necessary.UserID, []uuid.UUID{tag.AccountGroupID}); err != nil {
		return id, err
	}

	// Создаем подкатегорию
	return s.tagRepository.CreateTag(ctx, tag.ConvertToRepoReq())
}
