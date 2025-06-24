package service

import (
	"context"

	"server/internal/services/account/model"
)

// DeleteAccount удаляет счет
func (s *AccountService) DeleteAccount(ctx context.Context, req model.DeleteAccountReq) error {
	ctx, span := tracer.Start(ctx, "DeleteAccount")
	defer span.End()

	// Проверяем доступ пользователя к счету
	if err := s.CheckAccess(ctx, req.Necessary.UserID, []uint32{req.ID}); err != nil {
		return err
	}

	// Удаляем счет
	return s.accountRepository.DeleteAccount(ctx, req.ID)
}
