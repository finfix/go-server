package service

import (
	"context"

	"pkg/slices"
	"server/internal/utils/errors"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	accountModel "server/internal/modules/account/model"
	accountRepoModel "server/internal/modules/account/repository/model"
	auditLogModel "server/internal/modules/auditLog/model"
	transactionModel "server/internal/modules/transaction/model"
	"server/internal/modules/transaction/service/utils"

	"github.com/google/uuid"
)

// UpdateTransaction редактирует транзакцию
func (s *TransactionService) UpdateTransaction(ctx context.Context, fields transactionModel.UpdateTransactionReq) error {
	ctx, span := tracer.Start(ctx, "UpdateTransaction")
	defer span.End()

	// Проверяем доступ пользователя к транзакции
	if err := s.CheckAccess(ctx, fields.Necessary.UserID, []uuid.UUID{fields.ID}); err != nil {
		return err
	}

	// Получаем транзакцию
	transactions, err := s.transactionRepository.GetTransactions(ctx, transactionModel.GetTransactionsReq{ //nolint:exhaustruct
		IDs: []uuid.UUID{fields.ID},
	})
	if err != nil {
		return err
	}
	if len(transactions) == 0 {
		return errors.NotFound.New("Транзакция не найдена").
			WithContextParams(ctx).
			WithParams("ID", fields.ID)
	}
	transaction := transactions[0]

	// Сохраняем слепок "до" отдельно, так как ниже transaction мутируется для валидации
	transactionBefore := transaction

	// Если в запросе есть изменение счетов, то проверяем доступ пользователя к ним
	if fields.AccountFromID != nil || fields.AccountToID != nil {
		if fields.AccountFromID != nil {
			transaction.AccountFromID = *fields.AccountFromID
		}
		if fields.AccountToID != nil {
			transaction.AccountToID = *fields.AccountToID
		}

		// Проверяем доступ пользователя к счетам
		if err = s.accountService.CheckAccess(ctx, fields.Necessary.UserID, []uuid.UUID{transaction.AccountFromID, transaction.AccountToID}); err != nil {
			return err
		}

		// Получаем счета
		_accounts, err := s.accountRepository.GetAccounts(ctx, accountRepoModel.GetAccountsReq{ //nolint:exhaustruct
			IDs: []uuid.UUID{transaction.AccountFromID, transaction.AccountToID},
		})
		if err != nil {
			return err
		}
		accountsMap := slices.ToMap(_accounts, func(account accountModel.Account) uuid.UUID { return account.ID })

		// Проверяем, что счета не являются родительскими: баланс родительского счета - это сумма
		// балансов дочерних счетов, а не собственные транзакции, поэтому менять его напрямую нельзя
		if accountsMap[transaction.AccountFromID].IsParent || accountsMap[transaction.AccountToID].IsParent {
			return errors.BadRequest.New("Нельзя создать транзакцию для родительского счета").
				WithContextParams(ctx).
				WithParams(
					"AccountFromID", transaction.AccountFromID,
					"AccountToID", transaction.AccountToID,
				)
		}

		// Проверяем, что счета находятся в одной группе
		if accountsMap[transaction.AccountFromID].AccountGroupID != accountsMap[transaction.AccountToID].AccountGroupID {
			return errors.BadRequest.New("Счета находятся в разных группах").
				WithContextParams(ctx).
				WithParams(
					"AccountFromID", transaction.AccountFromID,
					"AccountToID", transaction.AccountToID,
				)
		}

		// Проверяем соответствие типов счета и типа транзакции
		if err = utils.TransactionAndAccountTypesValidation(ctx,
			accountsMap[transaction.AccountFromID],
			accountsMap[transaction.AccountToID],
			transaction.Type,
		); err != nil {
			return err
		}

		// Определяем группу счетов из самих счетов, не доверяя предыдущему значению транзакции
		accountGroupID := accountsMap[transaction.AccountFromID].AccountGroupID
		fields.AccountGroupID = &accountGroupID
	}

	// Группа счетов, к которой относится транзакция после применения изменений
	accountGroupID := transactionBefore.AccountGroupID
	if fields.AccountGroupID != nil {
		accountGroupID = *fields.AccountGroupID
	}

	return s.generalRepository.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Если в запросе есть изменение тегов
		if fields.TagIDs != nil {
			if err := s.updateTransactionTags(ctxTx, fields.Necessary.UserID, fields.ID, accountGroupID, *fields.TagIDs); err != nil {
				return err
			}
		}

		// Изменяем данные транзакции
		if err := s.transactionRepository.UpdateTransaction(ctxTx, fields); err != nil {
			return err
		}

		// Получаем актуальную транзакцию из БД для слепка "после" в аудит-логе
		transactionAfter, err := slices.FirstWithError(s.transactionRepository.GetTransactions(ctxTx, transactionModel.GetTransactionsReq{ //nolint:exhaustruct
			IDs: []uuid.UUID{fields.ID},
		}))
		if err != nil {
			return err
		}

		// Фиксируем изменение транзакции в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:         auditLogEntity.Transaction,
			Method:         auditLogMethod.Update,
			EntityID:       fields.ID.String(),
			Before:         transactionBefore,
			After:          transactionAfter,
			UserID:         fields.Necessary.UserID,
			DeviceID:       fields.Necessary.DeviceID,
			AccountGroupID: &transactionAfter.AccountGroupID,
		})
	})
}
