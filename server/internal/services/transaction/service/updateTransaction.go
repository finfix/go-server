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

// UpdateTransaction редактирует транзакцию
func (s *TransactionService) UpdateTransaction(ctx context.Context, fields transactionModel.UpdateTransactionReq) error {
	ctx, span := tracer.Start(ctx, "UpdateTransaction")
	defer span.End()

	// Проверяем доступ пользователя к транзакции
	if err := s.CheckAccess(ctx, fields.Necessary.UserID, []uint32{fields.ID}); err != nil {
		return err
	}

	// Получаем транзакцию
	transactions, err := s.transactionRepository.GetTransactions(ctx, transactionModel.GetTransactionsReq{ //nolint:exhaustruct
		IDs: []uint32{fields.ID},
	})
	if err != nil {
		return err
	}
	if len(transactions) == 0 {
		return errors.NotFound.New("Транзакция не найдена",
			errors.ParamsOption("ID", fields.ID),
		)
	}
	transaction := transactions[0]

	// Если в запросе есть изменение счетов, то проверяем доступ пользователя к ним
	if fields.AccountFromID != nil || fields.AccountToID != nil {
		if fields.AccountFromID != nil {
			transaction.AccountFromID = *fields.AccountFromID
		}
		if fields.AccountToID != nil {
			transaction.AccountToID = *fields.AccountToID
		}

		// Проверяем доступ пользователя к счетам
		if err = s.accountService.CheckAccess(ctx, fields.Necessary.UserID, []uint32{transaction.AccountFromID, transaction.AccountToID}); err != nil {
			return err
		}

		// Получаем счета
		_accounts, err := s.accountRepository.GetAccounts(ctx, accountRepoModel.GetAccountsReq{ //nolint:exhaustruct
			IDs: []uint32{transaction.AccountFromID, transaction.AccountToID},
		})
		if err != nil {
			return err
		}
		accountsMap := slices.ToMap(_accounts, func(account accountModel.Account) uint32 { return account.ID })

		// Проверяем соответствие типов счета и типа транзакции
		if err = utils.TransactionAndAccountTypesValidation(
			accountsMap[transaction.AccountFromID],
			accountsMap[transaction.AccountToID],
			transaction.Type,
		); err != nil {
			return err
		}
	}

	return s.generalRepository.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Если в запросе есть изменение тегов
		if fields.TagIDs != nil {
			if err := s.updateTransactionTags(ctxTx, fields.Necessary.UserID, fields.ID, *fields.TagIDs); err != nil {
				return err
			}
		}

		// Изменяем данные транзакции
		return s.transactionRepository.UpdateTransaction(ctxTx, fields)
	})
}
