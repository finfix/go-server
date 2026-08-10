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

// CreateAccount создает новый счет
func (s *AccountService) CreateAccount(ctx context.Context, accountToCreate model.CreateAccountReq) (res model.CreateAccountRes, err error) {
	ctx, span := tracer.Start(ctx, "CreateAccount")
	defer span.End()

	// Проверяем доступ пользователя к группе счетов
	if err = s.accountGroupService.CheckAccess(ctx, accountToCreate.Necessary.UserID, []uuid.UUID{accountToCreate.AccountGroupID}); err != nil {
		return res, err
	}

	// Проверяем, можно ли привязать счет к родительскому счету
	if accountToCreate.ParentAccountID != nil {

		// Представляем, что счет уже создан
		account := accountToCreate.ConvertToAccount()

		// Проверяем возможность привязки
		if err = s.ValidateUpdateParentAccountID(ctx, account, *accountToCreate.ParentAccountID, accountToCreate.Necessary.UserID); err != nil {
			return res, err
		}
	}

	err = s.transactor.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Создаем счет
		if err := s.accountRepository.CreateAccount(ctxTx, accountToCreate.ConvertToRepoReq()); err != nil {
			return err
		}

		// Получаем созданный счет из БД для слепка "после" в аудит-логе
		accountAfter, err := slices.FirstWithError(s.accountRepository.GetAccounts(ctxTx, accountRepoModel.GetAccountsReq{IDs: []uuid.UUID{accountToCreate.ID}})) //nolint:exhaustruct
		if err != nil {
			return err
		}

		// Фиксируем создание счета в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.Account,
			Method:   auditLogMethod.Create,
			EntityID: accountToCreate.ID.String(),
			Before:   nil,
			After:    accountAfter,
			UserID:   accountToCreate.Necessary.UserID,
			DeviceID: accountToCreate.Necessary.DeviceID,
		})
	})
	if err != nil {
		return res, err
	}

	return res, nil
}
