package service

import (
	"context"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	"server/internal/modules/account/model"
	accountRepoModel "server/internal/modules/account/repository/model"
	auditLogModel "server/internal/modules/auditLog/model"

	"github.com/google/uuid"
)

// DeleteAccount удаляет счет
func (s *AccountService) DeleteAccount(ctx context.Context, req model.DeleteAccountReq) error {
	ctx, span := tracer.Start(ctx, "DeleteAccount")
	defer span.End()

	// Проверяем доступ пользователя к счету
	if err := s.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{req.ID}); err != nil {
		return err
	}

	// Получаем счет для слепка "до" в аудит-логе
	account, err := slices.FirstWithError(s.accountRepository.GetAccounts(ctx, accountRepoModel.GetAccountsReq{IDs: []uuid.UUID{req.ID}})) //nolint:exhaustruct
	if err != nil {
		return err
	}

	return s.transactor.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Удаляем счет
		if err := s.accountRepository.DeleteAccount(ctxTx, req.ID); err != nil {
			return err
		}

		// Фиксируем удаление счета в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.Account,
			Method:   auditLogMethod.Delete,
			EntityID: req.ID.String(),
			Before:   account,
			After:    nil,
			UserID:   req.Necessary.UserID,
			DeviceID: req.Necessary.DeviceID,
		})
	})
}
