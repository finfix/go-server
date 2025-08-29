package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *AccountService) CheckAccess(ctx context.Context, userID uuid.UUID, accountIDs []uuid.UUID) error {
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
