package service

import (
	"context"
	"server/internal/enum/accountType"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"server/internal/utils/errors"

	"server/internal/modules/account/model"
	accountRepoModel "server/internal/modules/account/repository/model"
)

func (s *AccountService) calculateRemainders(ctx context.Context, filters model.GetAccountsReq) (map[uuid.UUID]decimal.Decimal, error) {
	ctx, span := tracer.Start(ctx, "calculateRemainders")
	defer span.End()

	// Считаем балансы обычных и долговых счетов
	calculatedRemainders, err := s.accountRepository.GetSumAllTransactionsToAccount(ctx, accountRepoModel.CalculateRemaindersAccountsReq{ //nolint:exhaustruct
		AccountGroupIDs: filters.AccountGroupIDs,
		Types: []accountType.AccountType{
			accountType.Debt,
			accountType.Regular,
		},
	})
	if err != nil {
		return nil, err
	}

	// Если тип счета расход, доход или балансировочный (или типа нет), то обязательно должен быть указан интервал дат
	if filters.Type == nil || *filters.Type == accountType.Earnings || *filters.Type == accountType.Expense || *filters.Type == accountType.Balancing {

		if filters.DateFrom == nil || filters.DateTo == nil {
			return nil, errors.BadRequest.New("dateFrom and dateTo must be specified").WithContextParams(ctx)
		}

		// Считаем расходы и доходы за указанный период или даты
		earnAndExp, err := s.accountRepository.GetSumAllTransactionsToAccount(ctx, accountRepoModel.CalculateRemaindersAccountsReq{ //nolint:exhaustruct
			AccountGroupIDs: filters.AccountGroupIDs,
			Types: []accountType.AccountType{
				accountType.Earnings,
				accountType.Expense,
				accountType.Balancing,
			},
			DateFrom: filters.DateFrom,
			DateTo:   filters.DateTo,
		})
		if err != nil {
			return nil, err
		}

		// Добавляем балансы расходов, доходов и балансировочных счетов к остаткам
		for id, amount := range earnAndExp {
			calculatedRemainders[id] = amount
		}
	}
	return calculatedRemainders, err
}
