package service

import (
	"context"

	"go.opentelemetry.io/otel"

	accountGroupModel "server/internal/services/accountGroup/model"
	accountGroupRepository "server/internal/services/accountGroup/repository"
	accountGroupRepoModel "server/internal/services/accountGroup/repository/model"
	"server/internal/services/transactor"
	userService "server/internal/services/user/service"
)

var tracer = otel.Tracer("/server/internal/services/accountGroup/service")

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
}

var _ AccountGroupRepository = new(accountGroupRepository.AccountGroupRepository)

type AccountGroupRepository interface {
	CreateAccountGroup(context.Context, accountGroupRepoModel.CreateAccountGroupReq) (uint32, uint32, error)
	GetAccountGroups(context.Context, accountGroupModel.GetAccountGroupsReq) ([]accountGroupModel.AccountGroup, error)
	UpdateAccountGroup(context.Context, accountGroupModel.UpdateAccountGroupReq) error
	DeleteAccountGroup(ctx context.Context, id uint32) error

	LinkUserToAccountGroup(ctx context.Context, userID, accountGroupID uint32) error
	UnlinkUserFromAccountGroup(ctx context.Context, userID, accountGroupID uint32) error
}

var _ UserService = new(userService.UserService)

type UserService interface {
	GetAccessedAccountGroups(ctx context.Context, userID uint32) (accesses []uint32, err error)
}

type AccountGroupService struct {
	userService            UserService
	accountGroupRepository AccountGroupRepository
	transactor             Transactor
}

func NewAccountGroupService(
	accountGroup AccountGroupRepository,
	transactor Transactor,
	userService UserService,
) *AccountGroupService {
	return &AccountGroupService{
		accountGroupRepository: accountGroup,
		transactor:             transactor,
		userService:            userService,
	}
}
