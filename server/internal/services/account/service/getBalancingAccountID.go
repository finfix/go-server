package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"pkg/pointer"
	"server/internal/utils/errors"

	"server/internal/services/account/model"
	"server/internal/services/account/model/accountType"
	accountRepoModel "server/internal/services/account/repository/model"
)

// GetBalancingAccountID получает ID балансировочного счета, подходящего для конкретного счета
func (s *AccountService) GetBalancingAccountID(ctx context.Context, account model.Account) (balancingAccountID uuid.UUID, serialNumber uint32, wasCreate bool, err error) {
	ctx, span := tracer.Start(ctx, "GetBalancingAccountID")
	defer span.End()

	// Получаем балансировочный счет группы в нужной валюте, чтобы создать для нее транзакцию
	balancingAccounts, err := s.accountRepository.GetAccounts(ctx, accountRepoModel.GetAccountsReq{ //nolint:exhaustruct
		Types:           []accountType.Type{accountType.Balancing},
		AccountGroupIDs: []uuid.UUID{account.AccountGroupID},
		Currencies:      []string{account.Currency},
		IsParent:        pointer.Pointer(false),
	})
	if err != nil {
		return balancingAccountID, serialNumber, wasCreate, err
	}
	// Если счет найден
	if len(balancingAccounts) != 0 {
		// Возвращаем его ID
		return balancingAccounts[0].ID, serialNumber, wasCreate, nil
	}

	// Получаем общий балансировочный счет
	parentBalancingAccounts, err := s.accountRepository.GetAccounts(ctx, accountRepoModel.GetAccountsReq{ //nolint:exhaustruct
		Types:           []accountType.Type{accountType.Balancing},
		AccountGroupIDs: []uuid.UUID{account.AccountGroupID},
		IsParent:        pointer.Pointer(true),
	})
	if err != nil {
		return balancingAccountID, serialNumber, wasCreate, err
	}

	var parentBalancingAccount model.Account

	// Если общий балансировочный счет не найден
	if len(parentBalancingAccounts) == 0 {
		return balancingAccountID, serialNumber, wasCreate, errors.InternalServer.New("Родительский балансировочный счет не найден").
			WithContextParams(ctx).
			WithParams(
				"accountID", account,
				"accountGroupID", account.AccountGroupID,
			)
	}

	parentBalancingAccount = parentBalancingAccounts[0]

	// Создаем балансировочный счет
	balancingAccountID, serialNumber, err = s.accountRepository.CreateAccount(ctx, accountRepoModel.CreateAccountReq{
		Budget: accountRepoModel.CreateReqBudget{
			Amount:         decimal.Zero,
			GradualFilling: false,
			FixedSum:       decimal.Zero,
			DaysOffset:     0,
		},
		Name:               "Балансировочный",
		Visible:            parentBalancingAccount.Visible,
		IconID:             parentBalancingAccount.IconID,
		Type:               accountType.Balancing,
		Currency:           account.Currency,
		AccountGroupID:     parentBalancingAccount.AccountGroupID,
		AccountingInHeader: parentBalancingAccount.AccountingInHeader,
		AccountingInCharts: true,
		IsParent:           false,
		ParentAccountID:    &parentBalancingAccount.ID,
		UserID:             account.CreatedByUserID,
		DatetimeCreate:     time.Now(),
	})
	if err != nil {
		return balancingAccountID, serialNumber, wasCreate, err
	}
	wasCreate = true
	return balancingAccountID, serialNumber, wasCreate, nil
}
