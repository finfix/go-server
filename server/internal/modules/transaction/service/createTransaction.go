package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/slices"
	"server/internal/utils/errors"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	accountModel "server/internal/modules/account/model"
	accountRepoModel "server/internal/modules/account/repository/model"
	auditLogModel "server/internal/modules/auditLog/model"
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

	// Проверяем, что счета не являются родительскими: баланс родительского счета - это сумма
	// балансов дочерних счетов, а не собственные транзакции, поэтому менять его напрямую нельзя
	if accountsMap[transaction.AccountFromID].IsParent || accountsMap[transaction.AccountToID].IsParent {
		return transactionModel.CreateTransactionRes{}, errors.BadRequest.New("Нельзя создать транзакцию для родительского счета").
			WithContextParams(ctx).
			WithParams(
				"AccountFromID", transaction.AccountFromID,
				"AccountToID", transaction.AccountToID,
			)
	}

	// Проверяем, может ли пользователь использовать счета
	if err = utils.TransactionAndAccountTypesValidation(ctx,
		accountsMap[transaction.AccountFromID],
		accountsMap[transaction.AccountToID],
		transaction.Type,
	); err != nil {
		return transactionModel.CreateTransactionRes{}, err
	}

	// Получаем разрешения счетов
	permissionsArr, err := accountModel.GetAccountsPermissions(accountsMap[transaction.AccountFromID], accountsMap[transaction.AccountToID])
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

	// Определяем группу счетов из самих счетов, не доверяя значению из запроса
	transaction.AccountGroupID = accountsMap[transaction.AccountFromID].AccountGroupID

	err = s.generalRepository.WithSyncGate(ctx, transaction.Necessary.UserID, transaction.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Создаем транзакцию
		id, err := s.transactionRepository.CreateTransaction(ctxTx, transaction.ConvertToRepoReq())
		if err != nil {
			return 0, err
		}

		// Если переданы теги
		if len(transaction.TagIDs) != 0 {
			if err = s.updateTransactionTags(ctxTx, transaction.Necessary.UserID, id, transaction.AccountGroupID, transaction.TagIDs); err != nil {
				return 0, err
			}
		}

		// Получаем созданную транзакцию из БД для слепка "после" в аудит-логе
		transactionAfter, err := slices.FirstWithError(s.transactionRepository.GetTransactions(ctxTx, transactionModel.GetTransactionsReq{ //nolint:exhaustruct
			IDs: []uuid.UUID{id},
		}))
		if err != nil {
			return 0, err
		}

		// Фиксируем создание транзакции в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:         auditLogEntity.Transaction,
			Method:         auditLogMethod.Create,
			EntityID:       id.String(),
			Before:         nil,
			After:          transactionAfter,
			UserID:         transaction.Necessary.UserID,
			DeviceID:       transaction.Necessary.DeviceID,
			AccountGroupID: &transactionAfter.AccountGroupID,
		})
	})
	if err != nil {
		return transactionModel.CreateTransactionRes{}, err
	}

	return transactionModel.CreateTransactionRes{}, nil
}
