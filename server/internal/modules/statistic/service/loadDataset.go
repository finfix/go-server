package service

import (
	"context"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	accountModel "server/internal/modules/account/model"
	settingsModel "server/internal/modules/settings/model"
	tagModel "server/internal/modules/tag/model"
	transactionModel "server/internal/modules/transaction/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

// loadDataset собирает счета, транзакции, теги и курсы выбранных групп счетов и определяет валюту
// ответа: явно переданную, иначе валюту единственной группы, иначе валюту пользователя
func (s *StatisticService) loadDataset(ctx context.Context, userInformation necessary.NecessaryUserInformation, accountGroupIDs []uuid.UUID, currency string) (data dataset, err error) {
	ctx, span := tracer.Start(ctx, "loadDataset")
	defer span.End()

	// Группы счетов: переданные проверяем на доступ, иначе берем все доступные
	if len(accountGroupIDs) != 0 {
		if err = s.accountGroupService.CheckAccess(ctx, userInformation.UserID, accountGroupIDs); err != nil {
			return data, err
		}
	} else {
		if accountGroupIDs, err = s.userToAccountGroupService.GetAccessedAccountGroups(ctx, userInformation.UserID); err != nil {
			return data, err
		}
	}
	data.accountGroupIDs = accountGroupIDs
	data.accounts = make(map[uuid.UUID]accountModel.Account)
	data.tagsByTransaction = make(map[uuid.UUID][]uuid.UUID)
	data.rateHistory = make(map[string][]ratePoint)
	data.approximate = make(map[string]bool)

	// Текущие курсы валют
	currencies, err := s.settingsService.GetCurrencies(ctx)
	if err != nil {
		return data, err
	}
	data.rates = make(map[string]decimal.Decimal, len(currencies.Currencies))
	for _, item := range currencies.Currencies {
		data.rates[item.Slug] = item.Rate
	}

	// Валюта ответа
	if data.currency, err = s.resolveCurrency(ctx, userInformation, accountGroupIDs, currency); err != nil {
		return data, err
	}
	if _, ok := data.rates[data.currency]; !ok {
		return data, errors.BadRequest.New("unknown currency").WithParams("currency", data.currency).
			WithCustomHumanText("Неизвестная валюта")
	}
	if len(accountGroupIDs) == 0 {
		return data, nil
	}

	// Счета групп. Период нужен сервису счетов только для остатков расходов и доходов, здесь не используется
	dateFrom, dateTo := time.Unix(0, 0).UTC(), time.Now().UTC()
	accounts, err := s.accountService.GetAccounts(ctx, accountModel.GetAccountsReq{ //nolint:exhaustruct
		Necessary:       userInformation,
		AccountGroupIDs: accountGroupIDs,
		DateFrom:        &dateFrom,
		DateTo:          &dateTo,
	})
	if err != nil {
		return data, err
	}
	for _, account := range accounts {
		data.accounts[account.ID] = account
	}

	// История курсов валют счетов и валюты ответа
	currencyCodes := []string{data.currency}
	for _, account := range data.accounts {
		if !slices.Contains(currencyCodes, account.Currency) {
			currencyCodes = append(currencyCodes, account.Currency)
		}
	}
	history, err := s.settingsService.GetCurrencyRates(ctx, settingsModel.GetCurrencyRatesReq{Currencies: currencyCodes}) //nolint:exhaustruct
	if err != nil {
		return data, err
	}
	for _, point := range history {
		data.rateHistory[point.Currency] = append(data.rateHistory[point.Currency], ratePoint{at: point.ValidAt, rate: point.Rate})
	}

	// Транзакции всех доступных групп, оставляем выбранные
	transactions, err := s.transactionService.GetTransactions(ctx, transactionModel.GetTransactionsReq{ //nolint:exhaustruct
		Necessary: userInformation,
	})
	if err != nil {
		return data, err
	}
	for _, transaction := range transactions {
		if slices.Contains(accountGroupIDs, transaction.AccountGroupID) {
			data.transactions = append(data.transactions, transaction)
		}
	}

	// Теги транзакций
	tagsToTransactions, err := s.tagService.GetTagsToTransactions(ctx, tagModel.GetTagsToTransactionsReq{ //nolint:exhaustruct
		Necessary: userInformation,
	})
	if err != nil {
		return data, err
	}
	for _, link := range tagsToTransactions {
		data.tagsByTransaction[link.TransactionID] = append(data.tagsByTransaction[link.TransactionID], link.TagID)
	}

	return data, nil
}
