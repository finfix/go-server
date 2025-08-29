package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/slices"
	"server/internal/utils/errors"

	accountModel "server/internal/modules/account/model"
	accountRepoModel "server/internal/modules/account/repository/model"
	transactionModel "server/internal/modules/transaction/model"
	"server/internal/modules/transaction/service/utils"
)

// CreateTransaction создает новую транзакцию
func (s *TransactionService) CreateTransaction(ctx context.Context, transaction transactionModel.CreateTransactionReq) (transactionModel.CreateTransactionRes, error) {
	ctx, span := tracer.Start(ctx, "CreateTransaction")
	defer span.End()

	// Проверяем доступ пользователя к счетам
	if err := s.accountService.CheckAccess(ctx, transaction.Necessary.UserID, []uuid.UUID{transaction.AccountFromID, transaction.AccountToID}); err != nil {
		return transactionModel.CreateTransactionRes{}, err
	}

	// Получаем счета
	_accounts, err := s.accountRepository.GetAccounts(ctx, accountRepoModel.GetAccountsReq{ //nolint:exhaustruct
		IDs: []uuid.UUID{transaction.AccountFromID, transaction.AccountToID},
	})
	if err != nil {
		return transactionModel.CreateTransactionRes{}, err
	}
	accountsMap := slices.ToMap(_accounts, func(account accountModel.Account) uuid.UUID { return account.ID })

	// Проверяем, может ли пользователь использовать счета
	if err = utils.TransactionAndAccountTypesValidation(ctx,
		accountsMap[transaction.AccountFromID],
		accountsMap[transaction.AccountToID],
		transaction.Type,
	); err != nil {
		return transactionModel.CreateTransactionRes{}, err
	}

	// Получаем разрешения счетов
	permissionsArr, err := s.permissionsService.GetAccountsPermissions(ctx, accountsMap[transaction.AccountFromID], accountsMap[transaction.AccountToID])
	if err != nil {
		return transactionModel.CreateTransactionRes{}, err
	}

	// Проверяем, что счета можно использовать для создания транзакции
	if !permissionsArr[0].CreateTransaction || !permissionsArr[1].CreateTransaction {
		return transactionModel.CreateTransactionRes{}, errors.BadRequest.New("Нельзя создать транзакцию для этих счетов").
			WithContextParams(ctx).
			WithParams(
				"AccountFromID", transaction.AccountFromID,
				"AccountGroupFromID", accountsMap[transaction.AccountFromID].AccountGroupID,
				"AccountToID", transaction.AccountToID,
				"AccountGroupToID", accountsMap[transaction.AccountToID].AccountGroupID,
			)
	}

	// Проверяем, что счета находятся в одной группе
	if accountsMap[transaction.AccountFromID].AccountGroupID != accountsMap[transaction.AccountToID].AccountGroupID {
		return transactionModel.CreateTransactionRes{}, errors.BadRequest.New("Счета находятся в разных группах").
			WithContextParams(ctx).
			WithParams(
				"AccountFromID", transaction.AccountFromID,
				"AccountToID", transaction.AccountToID,
			)
	}

	err = s.generalRepository.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Создаем транзакцию
		id, err := s.transactionRepository.CreateTransaction(ctx, transaction.ConvertToRepoReq())
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
	if err != nil {
		return transactionModel.CreateTransactionRes{}, err
	}

	return transactionModel.CreateTransactionRes{}, nil
}
