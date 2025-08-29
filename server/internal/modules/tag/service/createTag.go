package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/modules/tag/model"
)

// CreateTag создает новую подкатегорию
func (s *TagService) CreateTag(ctx context.Context, tag model.CreateTagReq) (model.CreateTagRes, error) {
	ctx, span := tracer.Start(ctx, "CreateTag")
	defer span.End()

	// Проверяем доступ пользователя к группам счетов
	if err := s.accountGroupService.CheckAccess(ctx, tag.Necessary.UserID, []uuid.UUID{tag.AccountGroupID}); err != nil {
		return model.CreateTagRes{}, err
	}

	// Создаем подкатегорию
	err := s.tagRepository.CreateTag(ctx, tag.ConvertToRepoReq())
	if err != nil {
		return model.CreateTagRes{}, err
	}

	return model.CreateTagRes{}, nil
}
