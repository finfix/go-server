package service

import (
	"context"

	"pkg/errors"
	"pkg/slices"

	accountModel "server/internal/services/account/model"
	accountRepoModel "server/internal/services/account/repository/model"
	transactionModel "server/internal/services/transaction/model"
	"server/internal/services/transaction/service/utils"
)

// CreateTransaction создает новую транзакцию
func (s *TransactionService) CreateTransaction(ctx context.Context, transaction transactionModel.CreateTransactionReq) (id uint32, err error) {
	ctx, span := tracer.Start(ctx, "CreateTransaction")
	defer span.End()

	// Проверяем доступ пользователя к счетам
	if err = s.accountService.CheckAccess(ctx, transaction.Necessary.UserID, []uint32{transaction.AccountFromID, transaction.AccountToID}); err != nil {
		return id, err
	}

	// Получаем счета
	_accounts, err := s.accountRepository.GetAccounts(ctx, accountRepoModel.GetAccountsReq{ //nolint:exhaustruct
		IDs: []uint32{transaction.AccountFromID, transaction.AccountToID},
	})
	if err != nil {
		return id, err
	}
	accountsMap := slices.ToMap(_accounts, func(account accountModel.Account) uint32 { return account.ID })

	// Проверяем, может ли пользователь использовать счета
	if err = utils.TransactionAndAccountTypesValidation(
		accountsMap[transaction.AccountFromID],
		accountsMap[transaction.AccountToID],
		transaction.Type,
	); err != nil {
		return id, err
	}

	// Получаем разрешения счетов
	permissionsArr, err := s.permissionsService.GetAccountsPermissions(ctx, accountsMap[transaction.AccountFromID], accountsMap[transaction.AccountToID])
	if err != nil {
		return id, err
	}

	// Проверяем, что счета можно использовать для создания транзакции
	if !permissionsArr[0].CreateTransaction || !permissionsArr[1].CreateTransaction {
		return id, errors.BadRequest.New("Нельзя создать транзакцию для этих счетов",
			errors.ParamsOption(
				"AccountFromID", transaction.AccountFromID,
				"AccountGroupFromID", accountsMap[transaction.AccountFromID].AccountGroupID,
				"AccountToID", transaction.AccountToID,
				"AccountGroupToID", accountsMap[transaction.AccountToID].AccountGroupID,
			),
		)
	}

	// Проверяем, что счета находятся в одной группе
	if accountsMap[transaction.AccountFromID].AccountGroupID != accountsMap[transaction.AccountToID].AccountGroupID {
		return id, errors.BadRequest.New("Счета находятся в разных группах",
			errors.ParamsOption(
				"AccountFromID", transaction.AccountFromID,
				"AccountToID", transaction.AccountToID,
			))
	}

	return id, s.generalRepository.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Создаем транзакцию
		id, err = s.transactionRepository.CreateTransaction(ctx, transaction.ConvertToRepoReq())
		if err != nil {
			return err
		}

		// Если переданы теги
		if len(transaction.TagIDs) != 0 {
			if err = s.updateTransactionTags(ctx, transaction.Necessary.UserID, id, transaction.TagIDs); err != nil {
				return err
			}
		}
		return nil
	})
}
