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

// CreateAccountGroup создает новую группу счетов
func (s *AccountGroupService) CreateAccountGroup(ctx context.Context, accountGroup model.CreateAccountGroupReq) (res model.CreateAccountGroupRes, err error) {
	ctx, span := tracer.Start(ctx, "CreateAccountGroup")
	defer span.End()

	// Создаем SQL-транзакцию
	return res, s.transactor.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Создаем счет
		serialNumber, err := s.accountGroupRepository.CreateAccountGroup(ctxTx, accountGroup.ConvertToRepoReq())
		if err != nil {
			return err
		}
		res.SerialNumber = serialNumber

		if err = s.accountGroupRepository.LinkUserToAccountGroup(ctxTx, accountGroup.Necessary.UserID, accountGroup.ID); err != nil {
			return err
		}

		// Получаем созданную группу счетов из БД для слепка "после" в аудит-логе
		accountGroupAfter, err := slices.FirstWithError(s.accountGroupRepository.GetAccountGroups(ctxTx, model.GetAccountGroupsReq{ //nolint:exhaustruct
			AccountGroupIDs: []uuid.UUID{accountGroup.ID},
		}))
		if err != nil {
			return err
		}

		// Фиксируем создание группы счетов в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.AccountGroup,
			Method:   auditLogMethod.Create,
			EntityID: accountGroup.ID.String(),
			Before:   nil,
			After:    accountGroupAfter,
			UserID:   accountGroup.Necessary.UserID,
			DeviceID: accountGroup.Necessary.DeviceID,
		})
	})
}
