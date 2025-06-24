package service

import "context"

func (s *TagService) CheckAccess(ctx context.Context, userID uint32, tagIDs []uint32) error {
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
