package service

import (
	"context"

	"pkg/slices"

	"server/internal/modules/account/model"
	accountRepoModel "server/internal/modules/account/repository/model"

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

	// Создаем SQL-транзакцию
	err = s.transactor.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Создаем счет
		serialNumber, err := s.accountRepository.CreateAccount(ctx, accountToCreate.ConvertToRepoReq())
		if err != nil {
			return err
		}
		res.SerialNumber = serialNumber

		// Если счет создался с остатком
		if !accountToCreate.Remainder.IsZero() {

			// Получаем счет
			account, err := slices.FirstWithError(s.accountRepository.GetAccounts(ctx,
				accountRepoModel.GetAccountsReq{ //nolint:exhaustruct
					IDs: []uuid.UUID{accountToCreate.ID},
				},
			))
			if err != nil {
				return err
			}

			// Меняем остаток счета созданием транзакции
			updateRes, err := s.ChangeAccountRemainder(ctxTx, account, accountToCreate.Remainder, accountToCreate.Necessary.UserID)
			if err != nil {
				return err
			}
			res.BalancingTransactionID = updateRes.BalancingTransactionID
			res.BalancingAccountID = updateRes.BalancingAccountID
			res.BalancingAccountSerialNumber = updateRes.BalancingAccountSerialNumber
		}

		return nil
	})
	if err != nil {
		return res, err
	}

	return res, nil
}
