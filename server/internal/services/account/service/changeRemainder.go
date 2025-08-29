package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"pkg/datetime"
	"server/internal/utils/errors"

	"server/internal/services/account/model"
	accountRepoModel "server/internal/services/account/repository/model"
	"server/internal/services/transaction/model/transactionType"
	transactionRepoModel "server/internal/services/transaction/repository/model"
)

func (s *AccountService) ChangeAccountRemainder(ctx context.Context, account model.Account, remainderToUpdate decimal.Decimal, userID uuid.UUID) (res model.UpdateAccountRes, err error) {
	ctx, span := tracer.Start(ctx, "ChangeAccountRemainder")
	defer span.End()

	// Получаем остаток счета
	remainders, err := s.accountRepository.GetSumAllTransactionsToAccount(ctx, accountRepoModel.CalculateRemaindersAccountsReq{
		IDs:             []uuid.UUID{account.ID},
		AccountGroupIDs: nil,
		Types:           nil,
		DateFrom:        nil,
		DateTo:          nil,
	})
	if err != nil {
		return res, err
	}

	// Проверяем, что остаток счета не равен написанному
	if remainderToUpdate == remainders[account.ID] {
		return res, errors.BadRequest.New("Остаток счета равен написанному").WithContextParams(ctx)
	}

	// Получаем балансировочный счет
	balancingAccountID, _, wasCreate, err := s.GetBalancingAccountID(ctx, account)
	if err != nil {
		return res, err
	}
	if wasCreate {
		res.BalancingAccountID = &balancingAccountID
		// res.BalancingTransactionID = &serialNumber
	}

	amount := remainderToUpdate.Sub(remainders[account.ID])

	// Создаем транзакцию балансировки
	balancingTransactionID, err := s.transactionRepository.CreateTransaction(ctx, transactionRepoModel.CreateTransactionReq{
		Type:               transactionType.Balancing,
		AmountFrom:         amount,
		AmountTo:           amount,
		Note:               "",
		AccountToID:        account.ID,
		AccountFromID:      balancingAccountID,
		DateTransaction:    datetime.Date{Time: time.Now()},
		IsExecuted:         true,
		CreatedByUserID:    userID,
		DatetimeCreate:     time.Now(),
		AccountingInCharts: true,
		AccountGroupID:     account.AccountGroupID,
	})
	if err != nil {
		return res, err
	}
	res.BalancingTransactionID = &balancingTransactionID

	return res, nil
}
