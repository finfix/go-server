package service

import (
	"context"

	"go.opentelemetry.io/otel"

	accountModel "server/internal/services/account/model"
	accountRepository "server/internal/services/account/repository"
	accountRepoModel "server/internal/services/account/repository/model"
	accountService "server/internal/services/account/service"
	"server/internal/services/accountPermissions/model"
	"server/internal/services/accountPermissions/service"
	tagModel "server/internal/services/tag/model"
	tagRepository "server/internal/services/tag/repository"
	tagService "server/internal/services/tag/service"
	transactionModel "server/internal/services/transaction/model"
	transactionRepository "server/internal/services/transaction/repository"
	transactionRepoModel "server/internal/services/transaction/repository/model"
	"server/internal/services/transactor"
	userService "server/internal/services/user/service"

	"github.com/google/uuid"
)

var tracer = otel.Tracer("/server/internal/services/transaction/service")

type TransactionService struct {
	transactionRepository TransactionRepository
	accountRepository     AccountRepository
	accountService        AccountService
	generalRepository     Transactor
	permissionsService    AccountPermissionsService
	tagRepository         TagRepository
	userService           UserService
	tagService            TagService
}

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
}

var _ TransactionRepository = new(transactionRepository.TransactionRepository)

type TransactionRepository interface {
	CreateTransaction(context.Context, transactionRepoModel.CreateTransactionReq) (uuid.UUID, error)
	UpdateTransaction(context.Context, transactionModel.UpdateTransactionReq) error
	DeleteTransaction(ctx context.Context, id, userID uuid.UUID) error
	GetTransactions(context.Context, transactionModel.GetTransactionsReq) (res []transactionModel.Transaction, err error)

	CheckAccess(ctx context.Context, accountGroupIDs, transactionIDs []uuid.UUID) error
}

var _ AccountPermissionsService = new(service.AccountPermissionsService)

type AccountPermissionsService interface {
	GetAccountsPermissions(context.Context, ...accountModel.Account) ([]model.AccountPermissions, error)
}

var _ AccountRepository = new(accountRepository.AccountRepository)

type AccountRepository interface {
	GetAccounts(context.Context, accountRepoModel.GetAccountsReq) ([]accountModel.Account, error)
}

var _ TagRepository = new(tagRepository.TagRepository)

type TagRepository interface {
	GetTagsToTransactions(context.Context, tagModel.GetTagsToTransactionsReq) ([]tagModel.TagToTransaction, error)
	LinkTagsToTransaction(context.Context, []uuid.UUID, uuid.UUID) error
	UnlinkTagsFromTransaction(context.Context, []uuid.UUID, uuid.UUID) error
}

var _ UserService = new(userService.UserService)

type UserService interface {
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

func NewTransactionService(
	transactionRepository TransactionRepository,
	accountRepository AccountRepository,
	transactor Transactor,
	accountPermissions AccountPermissionsService,
	tagRepository TagRepository,
	userService UserService,
	accountService AccountService,
	tagService TagService,
) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		accountRepository:     accountRepository,
		generalRepository:     transactor,
		permissionsService:    accountPermissions,
		tagRepository:         tagRepository,
		userService:           userService,
		accountService:        accountService,
		tagService:            tagService,
	}
}
