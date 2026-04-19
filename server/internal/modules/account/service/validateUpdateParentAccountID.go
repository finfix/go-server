package service

import (
	"context"

	"server/internal/utils/errors"

	"server/internal/modules/account/model"
	accountRepoModel "server/internal/modules/account/repository/model"

	"github.com/google/uuid"
)

func (s *AccountService) ValidateUpdateParentAccountID(ctx context.Context, account model.Account, parentAccountID, userID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "ValidateUpdateParentAccountID")
	defer span.End()

	if account.IsParent {
		return errors.BadRequest.New("Счет уже является родительским").
			WithContextParams(ctx).
			WithParams("accountID", account.ID)
	}

	if err := s.CheckAccess(ctx, userID, []uuid.UUID{parentAccountID}); err != nil {
		return err
	}

	// Получаем родительский счет
	parentAccounts, err := s.accountRepository.GetAccounts(ctx, accountRepoModel.GetAccountsReq{IDs: []uuid.UUID{parentAccountID}}) //nolint:exhaustruct
	if err != nil {
		return err
	}
	if len(parentAccounts) == 0 {
		return errors.NotFound.New("Родительский счет не найден").
			WithContextParams(ctx).
			WithParams("accountID", parentAccountID)
	}
	parentAccount := parentAccounts[0]

	// Проверяем, что указанный счет является родительским
	if parentAccount.ID != parentAccountID {
		return errors.BadRequest.New("Указанный счет не является родительским").
			WithContextParams(ctx).
			WithParams("accountID", parentAccountID)
	}

	// Проверяем, что счета находятся в одной группе
	if account.AccountGroupID != parentAccount.AccountGroupID {
		return errors.BadRequest.New("Счета находятся в разных группах").
			WithContextParams(ctx).
			WithParams(
				"childAccountID", account.ID,
				"childAccountGroupID", account.AccountGroupID,
				"parentAccountID", parentAccount.ID,
				"parentAccountGroupID", parentAccount.AccountGroupID,
			)
	}

	// Проверяем, что типы счетов совпадают
	if account.Type != parentAccount.Type {
		return errors.BadRequest.New("Типы счетов не совпадают").
			WithContextParams(ctx).
			WithParams(
				"childAccountID", account.ID,
				"childType", account.Type,
				"parentAccountID", parentAccount.ID,
				"parentType", parentAccount.Type,
			)
	}

	return nil
}
