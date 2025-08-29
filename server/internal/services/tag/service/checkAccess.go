package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *TagService) CheckAccess(ctx context.Context, userID uuid.UUID, tagIDs []uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "CheckAccess")
	defer span.End()

	// Получаем все доступные пользователю группы счетов
	accessedTagIDs, err := s.userService.GetAccessedAccountGroups(ctx, userID)
	if err != nil {
		return err
	}

	// Проверяем доступ пользователя к подкатегориям
	return s.tagRepository.CheckAccess(ctx, accessedTagIDs, tagIDs)
}
