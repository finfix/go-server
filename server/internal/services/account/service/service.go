package service

import (
	"context"

	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"

	accountModel "server/internal/services/account/model"
	accountRepository "server/internal/services/account/repository"
	accountRepoModel "server/internal/services/account/repository/model"
	accountGroupService "server/internal/services/accountGroup/service"
	accountPermissionsModel "server/internal/services/accountPermissions/model"
	accountPermisssionsService "server/internal/services/accountPermissions/service"
	transactionRepository "server/internal/services/transaction/repository"
	transactionRepoModel "server/internal/services/transaction/repository/model"
	"server/internal/services/transactor"
	userModel "server/internal/services/user/model"
	userRepository "server/internal/services/user/repository"
)

var tracer = otel.Tracer("/server/internal/services/account/service")

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
}

var _ AccountRepository = new(accountRepository.AccountRepository)

type AccountRepository interface {
	CreateAccount(context.Context, accountRepoModel.CreateAccountReq) (uint32, uint32, error)
	GetAccounts(context.Context, accountRepoModel.GetAccountsReq) ([]accountModel.Account, error)
	UpdateAccount(context.Context, map[uint32]accountRepoModel.UpdateAccountReq) error
	DeleteAccount(ctx context.Context, id uint32) error

	ChangeSerialNumbers(ctx context.Context, accountGroupID, oldValue, newValue uint32) error
	GetSumAllTransactionsToAccount(context.Context, accountRepoModel.CalculateRemaindersAccountsReq) (map[uint32]decimal.Decimal, error)

	CheckAccess(context.Context, []uint32, []uint32) error
}

var _ TransactionRepository = new(transactionRepository.TransactionRepository)

type TransactionRepository interface {
	CreateTransaction(context.Context, transactionRepoModel.CreateTransactionReq) (uint32, error)
}

var _ UserRepository = new(userRepository.UserRepository)

type UserRepository interface {
	GetUsers(context.Context, userModel.GetUsersReq) ([]userModel.User, error)
}

var _ AccountPermissionsService = new(accountPermisssionsService.AccountPermissionsService)

type AccountPermissionsService interface {
	GetAccountsPermissions(context.Context, ...accountModel.Account) ([]accountPermissionsModel.AccountPermissions, error)
}

var _ AccountGroupService = new(accountGroupService.AccountGroupService)

type AccountGroupService interface {
	CheckAccess(context.Context, uint32, []uint32) error
}

var _ UserService = new(userRepository.UserRepository)

type UserService interface {
	GetAccessedAccountGroups(ctx context.Context, userID uint32) (accesses []uint32, err error)
}

type AccountService struct {
	accountRepository         AccountRepository
	transactor                Transactor
	transactionRepository     TransactionRepository
	userRepository            UserRepository
	accountPermissionsService AccountPermissionsService
	accountGroupService       AccountGroupService
	userService               UserService
}

func NewAccountService(
	accountRepository AccountRepository,
	transactor Transactor,
	transactionRepository TransactionRepository,
	userRepository UserRepository,
	accountPermissionsService AccountPermissionsService,
	accountGroupsService AccountGroupService,
	userService UserService,
) *AccountService {
	return &AccountService{
		accountRepository:         accountRepository,
		transactor:                transactor,
		transactionRepository:     transactionRepository,
		userRepository:            userRepository,
		accountPermissionsService: accountPermissionsService,
		accountGroupService:       accountGroupsService,
		userService:               userService,
	}
}
