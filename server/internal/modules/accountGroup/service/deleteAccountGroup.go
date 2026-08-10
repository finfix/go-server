package service

import (
	"context"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	"server/internal/modules/accountGroup/model"
	auditLogModel "server/internal/modules/auditLog/model"

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

	// Получаем группу счетов для слепка "до" в аудит-логе
	accountGroup, err := slices.FirstWithError(s.accountGroupRepository.GetAccountGroups(ctx, model.GetAccountGroupsReq{ //nolint:exhaustruct
		AccountGroupIDs: []uuid.UUID{id.ID},
	}))
	if err != nil {
		return err
	}

	return s.transactor.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Отвязываем пользователя от группы счетов
		if err := s.accountGroupRepository.UnlinkUserFromAccountGroup(ctxTx, id.Necessary.UserID, id.ID); err != nil {
			return err
		}

		// Удаляем счет
		if err := s.accountGroupRepository.DeleteAccountGroup(ctxTx, id.ID); err != nil {
			return err
		}

		// Фиксируем удаление группы счетов в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.AccountGroup,
			Method:   auditLogMethod.Delete,
			EntityID: id.ID.String(),
			Before:   accountGroup,
			After:    nil,
			UserID:   id.Necessary.UserID,
			DeviceID: id.Necessary.DeviceID,
		})
	})
}
