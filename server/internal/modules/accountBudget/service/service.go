package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"

	accountModel "server/internal/modules/account/model"
	accountService "server/internal/modules/account/service"
	accountBudgetModel "server/internal/modules/accountBudget/model"
	"server/internal/modules/accountBudget/repository"
	accountGroupService "server/internal/modules/accountGroup/service"
	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
	userToAccountGroupService "server/internal/modules/userToAccountGroup/service"
)

var tracer = otel.Tracer("/server/internal/modules/accountBudget/service")

var _ AccountBudgetRepository = new(repository.AccountBudgetRepository)

// AccountBudgetRepository - интерфейс репозитория версий бюджета счетов
type AccountBudgetRepository interface {
	CreateAccountBudget(context.Context, accountBudgetModel.AccountBudget) error
	GetAccountBudgets(context.Context, repository.GetAccountBudgetsReq) ([]accountBudgetModel.AccountBudget, error)
}

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
	WithSyncGate(ctx context.Context, userID uuid.UUID, deviceID string, deviceSyncGate transactor.DeviceSyncGate, auditLogChangeChecker transactor.AuditLogChangeChecker, mutate func(ctxTx context.Context) (auditLogID uint32, err error)) error
}

var _ AccountService = new(accountService.AccountService)

// AccountService - интерфейс сервиса счетов
type AccountService interface {
	CheckAccess(ctx context.Context, userID uuid.UUID, accountIDs []uuid.UUID) error
	GetAccounts(context.Context, accountModel.GetAccountsReq) ([]accountModel.Account, error)
}

var _ AccountGroupService = new(accountGroupService.AccountGroupService)

// AccountGroupService - интерфейс сервиса групп счетов
type AccountGroupService interface {
	CheckAccess(ctx context.Context, userID uuid.UUID, accountGroupIDs []uuid.UUID) error
}

var _ UserToAccountGroupService = new(userToAccountGroupService.UserToAccountGroupService)

// UserToAccountGroupService - интерфейс сервиса связей пользователей с группами счетов
type UserToAccountGroupService interface {
	GetAccessedAccountGroups(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

var _ AuditLogService = new(auditLogService.AuditLogService)

// AuditLogService - интерфейс сервиса аудит-лога
type AuditLogService interface {
	TrackMutation(context.Context, auditLogModel.TrackMutationReq) (uint32, error)
	HasAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) (bool, error)
}

var _ UserService = new(userService.UserService)

// UserService - интерфейс сервиса пользователей, используется для отсечения мутаций от
// несинхронизированных устройств
type UserService interface {
	GetDeviceLastSyncedAuditLogIDForUpdate(ctx context.Context, userID uuid.UUID, deviceID string) (uint32, error)
	BumpDeviceCheckpoint(ctx context.Context, userID uuid.UUID, deviceID string, auditLogID uint32) error
}

// AccountBudgetService - сервис версий бюджета счетов
type AccountBudgetService struct {
	accountBudgetRepository   AccountBudgetRepository
	transactor                Transactor
	accountService            AccountService
	accountGroupService       AccountGroupService
	userToAccountGroupService UserToAccountGroupService
	auditLogService           AuditLogService
	userService               UserService
}

// NewAccountBudgetService создает новый сервис версий бюджета счетов
func NewAccountBudgetService(
	accountBudgetRepository AccountBudgetRepository,
	transactor Transactor,
	accountService AccountService,
	accountGroupService AccountGroupService,
	userToAccountGroupService UserToAccountGroupService,
	auditLogService AuditLogService,
	userService UserService,
) *AccountBudgetService {
	return &AccountBudgetService{
		accountBudgetRepository:   accountBudgetRepository,
		transactor:                transactor,
		accountService:            accountService,
		accountGroupService:       accountGroupService,
		userToAccountGroupService: userToAccountGroupService,
		auditLogService:           auditLogService,
		userService:               userService,
	}
}
