package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"

	accountModel "server/internal/modules/account/model"
	accountService "server/internal/modules/account/service"
	accountBudgetModel "server/internal/modules/accountBudget/model"
	accountBudgetService "server/internal/modules/accountBudget/service"
	accountGroupModel "server/internal/modules/accountGroup/model"
	accountGroupService "server/internal/modules/accountGroup/service"
	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	pendingLinkedTransferModel "server/internal/modules/pendingLinkedTransfer/model"
	pendingLinkedTransferService "server/internal/modules/pendingLinkedTransfer/service"
	settingsModel "server/internal/modules/settings/model"
	settingsService "server/internal/modules/settings/service"
	tagModel "server/internal/modules/tag/model"
	tagService "server/internal/modules/tag/service"
	transactionModel "server/internal/modules/transaction/model"
	transactionService "server/internal/modules/transaction/service"
	"server/internal/modules/transactor"
	userModel "server/internal/modules/user/model"
	userService "server/internal/modules/user/service"
)

var tracer = otel.Tracer("/server/internal/modules/sync/service")

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
}

var _ UserService = new(userService.UserService)

// UserService - интерфейс сервиса пользователей
type UserService interface {
	SetDevicePendingSync(ctx context.Context, userID uuid.UUID, deviceID string, pendingCheckpoint uint32, pendingSyncToken uuid.UUID) error
	ConfirmDeviceSync(ctx context.Context, userID uuid.UUID, deviceID string, pendingSyncToken uuid.UUID) error
	GetUsers(ctx context.Context, filters userModel.GetUsersReq) ([]userModel.User, error)
}

var _ AuditLogService = new(auditLogService.AuditLogService)

// AuditLogService - интерфейс сервиса аудит-лога
type AuditLogService interface {
	GetAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) ([]auditLogModel.AuditLog, error)
}

var _ TransactionService = new(transactionService.TransactionService)

// TransactionService - интерфейс сервиса транзакций
type TransactionService interface {
	GetTransactions(context.Context, transactionModel.GetTransactionsReq) ([]transactionModel.Transaction, error)
}

var _ AccountService = new(accountService.AccountService)

// AccountService - интерфейс сервиса счетов
type AccountService interface {
	GetAccounts(context.Context, accountModel.GetAccountsReq) ([]accountModel.Account, error)
}

var _ AccountGroupService = new(accountGroupService.AccountGroupService)

// AccountGroupService - интерфейс сервиса групп счетов
type AccountGroupService interface {
	GetAccountGroups(context.Context, accountGroupModel.GetAccountGroupsReq) ([]accountGroupModel.AccountGroup, error)
}

var _ TagService = new(tagService.TagService)

// TagService - интерфейс сервиса подкатегорий
type TagService interface {
	GetTags(context.Context, tagModel.GetTagsReq) ([]tagModel.Tag, error)
}

var _ AccountBudgetService = new(accountBudgetService.AccountBudgetService)

// AccountBudgetService - интерфейс сервиса версий бюджета счетов
type AccountBudgetService interface {
	GetAccountBudgetsByIDs(ctx context.Context, ids []uuid.UUID) ([]accountBudgetModel.AccountBudget, error)
}

var _ SettingsService = new(settingsService.SettingsService)

// SettingsService - интерфейс сервиса настроек (используется для получения валют)
type SettingsService interface {
	GetCurrencies(ctx context.Context) (settingsModel.GetCurrenciesRes, error)
}

var _ PendingLinkedTransferService = new(pendingLinkedTransferService.PendingLinkedTransferService)

// PendingLinkedTransferService - интерфейс сервиса переносов через счета-мосты
type PendingLinkedTransferService interface {
	GetPendingLinkedTransfers(context.Context, pendingLinkedTransferModel.GetPendingLinkedTransfersReq) ([]pendingLinkedTransferModel.PendingLinkedTransfer, error)
}

// SyncService - сервис синхронизации изменений между устройствами
type SyncService struct {
	transactor                   Transactor
	userService                  UserService
	auditLogService              AuditLogService
	transactionService           TransactionService
	accountService                AccountService
	accountGroupService           AccountGroupService
	tagService                    TagService
	accountBudgetService          AccountBudgetService
	settingsService                SettingsService
	pendingLinkedTransferService PendingLinkedTransferService
}

// NewSyncService создает новый сервис синхронизации
func NewSyncService(
	transactor Transactor,
	userService UserService,
	auditLogService AuditLogService,
	transactionService TransactionService,
	accountService AccountService,
	accountGroupService AccountGroupService,
	tagService TagService,
	accountBudgetService AccountBudgetService,
	settingsService SettingsService,
	pendingLinkedTransferService PendingLinkedTransferService,
) *SyncService {
	return &SyncService{
		transactor:                   transactor,
		userService:                  userService,
		auditLogService:              auditLogService,
		transactionService:           transactionService,
		accountService:               accountService,
		accountGroupService:          accountGroupService,
		tagService:                   tagService,
		accountBudgetService:         accountBudgetService,
		settingsService:              settingsService,
		pendingLinkedTransferService: pendingLinkedTransferService,
	}
}
