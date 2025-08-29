package service

import (
	"context"

	"server/internal/services/accountGroup/model"

	"github.com/google/uuid"
)

// UpdateAccountGroup обновляет группу счетов по конкретным полям
func (s *AccountGroupService) UpdateAccountGroup(ctx context.Context, updateReq model.UpdateAccountGroupReq) error {
	ctx, span := tracer.Start(ctx, "UpdateAccountGroup")
	defer span.End()

	// Проверяем доступ пользователя к группе счетов
	if err := s.CheckAccess(ctx, updateReq.Necessary.UserID, []uuid.UUID{updateReq.ID}); err != nil {
		return err
	}

	// Обновляем группу счетов
	return s.accountGroupRepository.UpdateAccountGroup(ctx, updateReq)
}
