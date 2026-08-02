package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	accountModel "server/internal/modules/account/model"
	accountRepository "server/internal/modules/account/repository"
	accountRepoModel "server/internal/modules/account/repository/model"
	accountGroupService "server/internal/modules/accountGroup/service"
	transactionRepository "server/internal/modules/transaction/repository"
	transactionRepoModel "server/internal/modules/transaction/repository/model"
	"server/internal/modules/transactor"
	userModel "server/internal/modules/user/model"
	userRepository "server/internal/modules/user/repository"
)

var tracer = otel.Tracer("/server/internal/modules/account/service")

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
}

var _ AccountRepository = new(accountRepository.AccountRepository)

type AccountRepository interface {
	CreateAccount(context.Context, accountRepoModel.CreateAccountReq) error
	GetAccounts(context.Context, accountRepoModel.GetAccountsReq) ([]accountModel.Account, error)
	UpdateAccount(context.Context, map[uuid.UUID]accountRepoModel.UpdateAccountReq) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error

	GetSumAllTransactionsToAccount(context.Context, accountRepoModel.CalculateRemaindersAccountsReq) (map[uuid.UUID]decimal.Decimal, error)

	CheckAccess(context.Context, []uuid.UUID, []uuid.UUID) error
}

var _ TransactionRepository = new(transactionRepository.TransactionRepository)

type TransactionRepository interface {
	CreateTransaction(context.Context, transactionRepoModel.CreateTransactionReq) (uuid.UUID, error)
}

var _ UserRepository = new(userRepository.UserRepository)

type UserRepository interface {
	GetUsers(context.Context, userModel.GetUsersReq) ([]userModel.User, error)
}

var _ AccountGroupService = new(accountGroupService.AccountGroupService)

type AccountGroupService interface {
	CheckAccess(context.Context, uuid.UUID, []uuid.UUID) error
}

var _ UserService = new(userRepository.UserRepository)

type UserService interface {
	GetAccessedAccountGroups(ctx context.Context, userID uuid.UUID) (accesses []uuid.UUID, err error)
}

type AccountService struct {
	accountRepository     AccountRepository
	transactor            Transactor
	transactionRepository TransactionRepository
	userRepository        UserRepository
	accountGroupService   AccountGroupService
	userService           UserService
}

func NewAccountService(
	accountRepository AccountRepository,
	transactor Transactor,
	transactionRepository TransactionRepository,
	userRepository UserRepository,
	accountGroupsService AccountGroupService,
	userService UserService,
) *AccountService {
	return &AccountService{
		accountRepository:     accountRepository,
		transactor:            transactor,
		transactionRepository: transactionRepository,
		userRepository:        userRepository,
		accountGroupService:   accountGroupsService,
		userService:           userService,
	}
}
