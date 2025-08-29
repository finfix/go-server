package service

import (
	"context"

	"server/internal/services/accountGroup/model"

	"github.com/google/uuid"
)

// DeleteAccountGroup удаляет группу счетов
func (s *AccountGroupService) DeleteAccountGroup(ctx context.Context, id model.DeleteAccountGroupReq) error {
	ctx, span := tracer.Start(ctx, "DeleteAccountGroup")
	defer span.End()

	// Проверяем доступ пользователя к счету
	if err := s.CheckAccess(ctx, id.Necessary.UserID, []uuid.UUID{id.ID}); err != nil {
		return err
	}

	// Отвязываем пользователя от группы счетов
	if err := s.accountGroupRepository.UnlinkUserFromAccountGroup(ctx, id.Necessary.UserID, id.ID); err != nil {
		return err
	}

	// Удаляем счет
	return s.accountGroupRepository.DeleteAccountGroup(ctx, id.ID)
}
