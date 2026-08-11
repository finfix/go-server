package service

import (
	"context"

	"go.opentelemetry.io/otel"

	accountModel "server/internal/modules/account/model"
	accountRepository "server/internal/modules/account/repository"
	accountRepoModel "server/internal/modules/account/repository/model"
	accountService "server/internal/modules/account/service"
	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	tagModel "server/internal/modules/tag/model"
	tagRepository "server/internal/modules/tag/repository"
	tagService "server/internal/modules/tag/service"
	transactionModel "server/internal/modules/transaction/model"
	transactionRepository "server/internal/modules/transaction/repository"
	transactionRepoModel "server/internal/modules/transaction/repository/model"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
	userToAccountGroupService "server/internal/modules/userToAccountGroup/service"

	"github.com/google/uuid"
)

var tracer = otel.Tracer("/server/internal/modules/transaction/service")

type TransactionService struct {
	transactionRepository     TransactionRepository
	accountRepository         AccountRepository
	accountService            AccountService
	generalRepository         Transactor
	tagRepository             TagRepository
	userToAccountGroupService UserToAccountGroupService
	tagService                TagService
	auditLogService           AuditLogService
	userService               UserService
}

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
	WithSyncGate(ctx context.Context, userID uuid.UUID, deviceID string, deviceSyncGate transactor.DeviceSyncGate, auditLogChangeChecker transactor.AuditLogChangeChecker, mutate func(ctxTx context.Context) (auditLogID uint32, err error)) error
}

var _ TransactionRepository = new(transactionRepository.TransactionRepository)

type TransactionRepository interface {
	CreateTransaction(context.Context, transactionRepoModel.CreateTransactionReq) (uuid.UUID, error)
	UpdateTransaction(context.Context, transactionModel.UpdateTransactionReq) error
	DeleteTransaction(ctx context.Context, id, userID uuid.UUID) error
	GetTransactions(context.Context, transactionModel.GetTransactionsReq) (res []transactionModel.Transaction, err error)

	CheckAccess(ctx context.Context, accountGroupIDs, transactionIDs []uuid.UUID) error
}

var _ AccountRepository = new(accountRepository.AccountRepository)

type AccountRepository interface {
	GetAccounts(context.Context, accountRepoModel.GetAccountsReq) ([]accountModel.Account, error)
}

var _ TagRepository = new(tagRepository.TagRepository)

type TagRepository interface {
	GetTags(context.Context, tagModel.GetTagsReq) ([]tagModel.Tag, error)
	GetTagsToTransactions(context.Context, tagModel.GetTagsToTransactionsReq) ([]tagModel.TagToTransaction, error)
	LinkTagsToTransaction(context.Context, []uuid.UUID, uuid.UUID) error
	UnlinkTagsFromTransaction(context.Context, []uuid.UUID, uuid.UUID) error
}

var _ UserToAccountGroupService = new(userToAccountGroupService.UserToAccountGroupService)

type UserToAccountGroupService interface {
	GetAccessedAccountGroups(ctx context.Context, userID uuid.UUID) (accesses []uuid.UUID, err error)
}

var _ AccountService = new(accountService.AccountService)

type AccountService interface {
	CheckAccess(ctx context.Context, userID uuid.UUID, accountIDs []uuid.UUID) error
}

var _ TagService = new(tagService.TagService)

type TagService interface {
	CheckAccess(ctx context.Context, userID uuid.UUID, tagIDs []uuid.UUID) error
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

func NewTransactionService(
	transactionRepository TransactionRepository,
	accountRepository AccountRepository,
	transactor Transactor,
	tagRepository TagRepository,
	userToAccountGroupService UserToAccountGroupService,
	accountService AccountService,
	tagService TagService,
	auditLogService AuditLogService,
	userService UserService,
) *TransactionService {
	return &TransactionService{
		transactionRepository:     transactionRepository,
		accountRepository:         accountRepository,
		generalRepository:         transactor,
		tagRepository:             tagRepository,
		userToAccountGroupService: userToAccountGroupService,
		accountService:            accountService,
		tagService:                tagService,
		auditLogService:           auditLogService,
		userService:               userService,
	}
}
