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

// UpdateAccountGroup обновляет группу счетов по конкретным полям
func (s *AccountGroupService) UpdateAccountGroup(ctx context.Context, updateReq model.UpdateAccountGroupReq) error {
	ctx, span := tracer.Start(ctx, "UpdateAccountGroup")
	defer span.End()

	// Проверяем доступ пользователя к группе счетов
	if err := s.CheckAccess(ctx, updateReq.Necessary.UserID, []uuid.UUID{updateReq.ID}); err != nil {
		return err
	}

	// Получаем группу счетов для слепка "до" в аудит-логе
	accountGroup, err := slices.FirstWithError(s.accountGroupRepository.GetAccountGroups(ctx, model.GetAccountGroupsReq{ //nolint:exhaustruct
		AccountGroupIDs: []uuid.UUID{updateReq.ID},
	}))
	if err != nil {
		return err
	}

	return s.transactor.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Обновляем группу счетов
		if err := s.accountGroupRepository.UpdateAccountGroup(ctxTx, updateReq); err != nil {
			return err
		}

		// Фиксируем изменение группы счетов в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.AccountGroup,
			Method:   auditLogMethod.Update,
			EntityID: updateReq.ID.String(),
			Before:   accountGroup,
			After:    updateReq,
			UserID:   updateReq.Necessary.UserID,
			DeviceID: updateReq.Necessary.DeviceID,
		})
	})
}
