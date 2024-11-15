package service

import "context"

func (s *AccountService) CheckAccess(ctx context.Context, userID uint32, accountIDs []uint32) error {
	ctx, span := tracer.Start(ctx, "CheckAccess")
	defer span.End()

	// Получаем доступные для пользователя группы счетов
	accessedAccountIDs, err := s.userService.GetAccessedAccountGroups(ctx, userID)
	if err != nil {
		return err
	}

	// Проверяем доступ пользователя к счетам
	return s.accountRepository.CheckAccess(ctx, accessedAccountIDs, accountIDs)
}
