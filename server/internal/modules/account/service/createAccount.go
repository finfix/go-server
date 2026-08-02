package service

import (
	"context"

	"server/internal/modules/account/model"

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

	// Создаем счет
	if err = s.accountRepository.CreateAccount(ctx, accountToCreate.ConvertToRepoReq()); err != nil {
		return res, err
	}

	return res, nil
}
