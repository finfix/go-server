// Package service — статистика доходов, расходов и наличности для MCP: те же разбивки и фильтры,
// что на экранах статистики и транзакций приложения. Данные берет у сервисов модулей через их
// интерфейсы, своего хранилища у модуля нет
package service

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	accountModel "server/internal/modules/account/model"
	accountService "server/internal/modules/account/service"
	accountBudgetModel "server/internal/modules/accountBudget/model"
	accountBudgetService "server/internal/modules/accountBudget/service"
	accountGroupModel "server/internal/modules/accountGroup/model"
	accountGroupService "server/internal/modules/accountGroup/service"
	settingsModel "server/internal/modules/settings/model"
	settingsService "server/internal/modules/settings/service"
	tagModel "server/internal/modules/tag/model"
	tagService "server/internal/modules/tag/service"
	transactionModel "server/internal/modules/transaction/model"
	transactionService "server/internal/modules/transaction/service"
	userModel "server/internal/modules/user/model"
	userService "server/internal/modules/user/service"
	userToAccountGroupService "server/internal/modules/userToAccountGroup/service"
)

var tracer = otel.Tracer("/server/internal/modules/statistic/service")

var _ TransactionService = new(transactionService.TransactionService)

// TransactionService — транзакции доступных пользователю групп счетов
type TransactionService interface {
	GetTransactions(context.Context, transactionModel.GetTransactionsReq) ([]transactionModel.Transaction, error)
}

var _ TagService = new(tagService.TagService)

// TagService — теги и их связи с транзакциями
type TagService interface {
	GetTagsToTransactions(context.Context, tagModel.GetTagsToTransactionsReq) ([]tagModel.TagToTransaction, error)
	GetTags(context.Context, tagModel.GetTagsReq) ([]tagModel.Tag, error)
}

var _ AccountService = new(accountService.AccountService)

// AccountService — счета с типами, валютами и иерархией
type AccountService interface {
	GetAccounts(context.Context, accountModel.GetAccountsReq) ([]accountModel.Account, error)
}

var _ AccountGroupService = new(accountGroupService.AccountGroupService)

// AccountGroupService — группы счетов и проверка доступа к ним
type AccountGroupService interface {
	CheckAccess(ctx context.Context, userID uuid.UUID, accountGroupIDs []uuid.UUID) error
	GetAccountGroups(context.Context, accountGroupModel.GetAccountGroupsReq) ([]accountGroupModel.AccountGroup, error)
}

var _ UserToAccountGroupService = new(userToAccountGroupService.UserToAccountGroupService)

// UserToAccountGroupService — группы счетов, доступные пользователю
type UserToAccountGroupService interface {
	GetAccessedAccountGroups(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

var _ SettingsService = new(settingsService.SettingsService)

// SettingsService — текущие курсы валют
type SettingsService interface {
	GetCurrencies(context.Context) (settingsModel.GetCurrenciesRes, error)
	GetCurrencyRates(context.Context, settingsModel.GetCurrencyRatesReq) ([]settingsModel.CurrencyRate, error)
}

var _ UserService = new(userService.UserService)

// UserService — валюта пользователя по умолчанию
type UserService interface {
	GetUser(context.Context, userModel.GetUserReq) (userModel.User, error)
}

var _ AccountBudgetService = new(accountBudgetService.AccountBudgetService)

// AccountBudgetService — версии бюджетов расходных счетов
type AccountBudgetService interface {
	GetAccountBudgets(context.Context, accountBudgetModel.GetAccountBudgetsReq) ([]accountBudgetModel.AccountBudget, error)
}

// StatisticService — расчет статистики
type StatisticService struct {
	transactionService        TransactionService
	tagService                TagService
	accountService            AccountService
	accountGroupService       AccountGroupService
	userToAccountGroupService UserToAccountGroupService
	settingsService           SettingsService
	userService               UserService
	accountBudgetService      AccountBudgetService
}

// NewStatisticService создает сервис статистики
func NewStatisticService(
	transactionService TransactionService,
	tagService TagService,
	accountService AccountService,
	accountGroupService AccountGroupService,
	userToAccountGroupService UserToAccountGroupService,
	settingsService SettingsService,
	userService UserService,
	accountBudgetService AccountBudgetService,
) *StatisticService {
	return &StatisticService{
		transactionService:        transactionService,
		tagService:                tagService,
		accountService:            accountService,
		accountGroupService:       accountGroupService,
		userToAccountGroupService: userToAccountGroupService,
		settingsService:           settingsService,
		userService:               userService,
		accountBudgetService:      accountBudgetService,
	}
}
