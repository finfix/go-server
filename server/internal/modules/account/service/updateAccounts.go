package service

import (
	"context"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	"server/internal/modules/account/model"
	accountRepoModel "server/internal/modules/account/repository/model"
	"server/internal/modules/account/service/utils"
	auditLogModel "server/internal/modules/auditLog/model"

	"github.com/google/uuid"
)

// UpdateAccount обновляет счет по конкретным полям
func (s *AccountService) UpdateAccount(ctx context.Context, updateReq model.UpdateAccountReq) (err error) {
	ctx, span := tracer.Start(ctx, "UpdateAccount")
	defer span.End()

	repoUpdateReqs := make(map[uuid.UUID]accountRepoModel.UpdateAccountReq)
	repoUpdateReqs[updateReq.ID] = updateReq.ConvertToRepoReq()

	// Проверяем доступ пользователя к счету
	if err = s.CheckAccess(ctx, updateReq.Necessary.UserID, []uuid.UUID{updateReq.ID}); err != nil {
		return err
	}

	// Получаем счет для слепка "до" в аудит-логе
	accountBefore, err := slices.FirstWithError(s.accountRepository.GetAccounts(ctx, accountRepoModel.GetAccountsReq{IDs: []uuid.UUID{updateReq.ID}})) //nolint:exhaustruct
	if err != nil {
		return err
	}

	// Проверяем, что входные данные не противоречат разрешениям
	permissions, err := model.GetAccountPermissions(accountBefore)
	if err != nil {
		return err
	}
	if err = utils.CheckAccountPermissionsForUpdate(ctx, updateReq, permissions); err != nil {
		return err
	}

	// Если привязываем счет к родительскому счету, проверяем возможность привязки
	if updateReq.ParentAccountID != nil && *updateReq.ParentAccountID != uuid.Nil {
		if err := s.ValidateUpdateParentAccountID(ctx, accountBefore, *updateReq.ParentAccountID, updateReq.Necessary.UserID); err != nil {
			return err
		}
	}

	return s.transactor.WithSyncGate(ctx, updateReq.Necessary.UserID, updateReq.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Обновляем счет
		if err := s.accountRepository.UpdateAccount(ctxTx, repoUpdateReqs); err != nil {
			return 0, err
		}

		// Получаем актуальный счет из БД для слепка "после" в аудит-логе
		accountAfter, err := slices.FirstWithError(s.accountRepository.GetAccounts(ctxTx, accountRepoModel.GetAccountsReq{IDs: []uuid.UUID{updateReq.ID}})) //nolint:exhaustruct
		if err != nil {
			return 0, err
		}

		// Фиксируем изменение счета в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:         auditLogEntity.Account,
			Method:         auditLogMethod.Update,
			EntityID:       updateReq.ID.String(),
			Before:         accountBefore,
			After:          accountAfter,
			UserID:         updateReq.Necessary.UserID,
			DeviceID:       updateReq.Necessary.DeviceID,
			AccountGroupID: &accountAfter.AccountGroupID,
		})
	})
}
