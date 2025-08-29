package service

import (
	"context"

	"pkg/slices"

	"server/internal/modules/account/model"
	accountRepoModel "server/internal/modules/account/repository/model"
	"server/internal/modules/account/service/utils"

	"github.com/google/uuid"
)

// UpdateAccount обновляет счет по конкретным полям
func (s *AccountService) UpdateAccount(ctx context.Context, updateReq model.UpdateAccountReq) (res model.UpdateAccountRes, err error) {
	ctx, span := tracer.Start(ctx, "UpdateAccount")
	defer span.End()

	repoUpdateReqs := make(map[uuid.UUID]accountRepoModel.UpdateAccountReq)
	repoUpdateReqs[updateReq.ID] = updateReq.ConvertToRepoReq()

	// Проверяем доступ пользователя к счету
	if err = s.CheckAccess(ctx, updateReq.Necessary.UserID, []uuid.UUID{updateReq.ID}); err != nil {
		return res, err
	}

	// Получаем счет
	account, err := slices.FirstWithError(s.accountRepository.GetAccounts(ctx, accountRepoModel.GetAccountsReq{IDs: []uuid.UUID{updateReq.ID}})) //nolint:exhaustruct
	if err != nil {
		return res, err
	}

	// Проверяем, что входные данные не противоречат разрешениям
	permissions, err := slices.FirstWithError(s.accountPermissionsService.GetAccountsPermissions(ctx, account))
	if err != nil {
		return res, err
	}
	if err = utils.CheckAccountPermissionsForUpdate(ctx, updateReq, permissions); err != nil {
		return res, err
	}

	// Проверяем, можно ли привязать счет к родительскому счету
	if updateReq.ParentAccountID != nil {

		// Если привязываем счет к родительскому счету
		if *updateReq.ParentAccountID != uuid.Nil {

			// Проверяем возможность привязки
			if err := s.ValidateUpdateParentAccountID(ctx, account, *updateReq.ParentAccountID, updateReq.Necessary.UserID); err != nil {
				return res, err
			}
			account.ParentAccountID = updateReq.ParentAccountID

		} else { // Если отвязываем счет от родительского счета
			account.ParentAccountID = nil
		}
	}

	// Получаем дочерние счета
	var childrenAccounts []model.Account
	if account.IsParent {
		childrenAccounts, err = s.accountRepository.GetAccounts(ctx, accountRepoModel.GetAccountsReq{ParentAccountIDs: []uuid.UUID{updateReq.ID}}) //nolint:exhaustruct
		if err != nil {
			return res, err
		}
	}

	// Получаем родительский счет
	var parentAccount *model.Account
	if account.ParentAccountID != nil {
		parentAccounts, err := s.accountRepository.GetAccounts(ctx, accountRepoModel.GetAccountsReq{IDs: []uuid.UUID{*account.ParentAccountID}}) //nolint:exhaustruct
		if err != nil {
			return res, err
		}
		parentAccount = &parentAccounts[0]
	}

	if updateReq.AccountingInHeader != nil {
		account.AccountingInHeader = *updateReq.AccountingInHeader
	}
	utils.HandleAccountingInHeader(
		repoUpdateReqs,
		account,
		childrenAccounts,
		parentAccount,
	)

	if updateReq.Visible != nil {
		account.Visible = *updateReq.Visible
	}
	utils.HandleVisible(
		repoUpdateReqs,
		account,
		childrenAccounts,
		parentAccount,
	)

	return res, s.transactor.WithinTransaction(ctx, func(ctxTx context.Context) error {
		res, err = s.updateAccounts(ctxTx, account, repoUpdateReqs, updateReq.Necessary.UserID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *AccountService) updateAccounts(ctx context.Context, account model.Account, updateReqs map[uuid.UUID]accountRepoModel.UpdateAccountReq, userID uuid.UUID) (res model.UpdateAccountRes, err error) {

	// Если передан остаток, редактируем его
	if updateReqs[account.ID].Remainder != nil {
		if res, err = s.ChangeAccountRemainder(
			ctx,
			account,
			*updateReqs[account.ID].Remainder,
			userID,
		); err != nil {
			return res, err
		}
	}

	// Если передан порядковый номер, то меняем порядковые номера остальных счетов
	if updateReqs[account.ID].SerialNumber != nil {
		if err = s.accountRepository.ChangeSerialNumbers(
			ctx,
			account.AccountGroupID,
			*updateReqs[account.ID].SerialNumber,
			account.SerialNumber,
		); err != nil {
			return res, err
		}
	}

	// Редактируем счет
	return res, s.accountRepository.UpdateAccount(ctx, updateReqs)
}
