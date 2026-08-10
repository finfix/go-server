package service

import (
	"context"

	"go.opentelemetry.io/otel"

	accountGroupModel "server/internal/modules/accountGroup/model"
	accountGroupRepository "server/internal/modules/accountGroup/repository"
	accountGroupRepoModel "server/internal/modules/accountGroup/repository/model"
	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"

	"github.com/google/uuid"
)

var tracer = otel.Tracer("/server/internal/modules/accountGroup/service")

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
}

var _ AccountGroupRepository = new(accountGroupRepository.AccountGroupRepository)

type AccountGroupRepository interface {
	CreateAccountGroup(context.Context, accountGroupRepoModel.CreateAccountGroupReq) (uint32, error)
	GetAccountGroups(context.Context, accountGroupModel.GetAccountGroupsReq) ([]accountGroupModel.AccountGroup, error)
	UpdateAccountGroup(context.Context, accountGroupModel.UpdateAccountGroupReq) error
	DeleteAccountGroup(ctx context.Context, id uuid.UUID) error

	LinkUserToAccountGroup(ctx context.Context, userID, accountGroupID uuid.UUID) error
	UnlinkUserFromAccountGroup(ctx context.Context, userID, accountGroupID uuid.UUID) error
}

var _ UserService = new(userService.UserService)

type UserService interface {
	GetAccessedAccountGroups(ctx context.Context, userID uuid.UUID) (accesses []uuid.UUID, err error)
}

var _ AuditLogService = new(auditLogService.AuditLogService)

// AuditLogService - интерфейс сервиса аудит-лога
type AuditLogService interface {
	TrackMutation(context.Context, auditLogModel.TrackMutationReq) error
}

type AccountGroupService struct {
	userService            UserService
	accountGroupRepository AccountGroupRepository
	transactor             Transactor
	auditLogService        AuditLogService
}

func NewAccountGroupService(
	accountGroup AccountGroupRepository,
	transactor Transactor,
	userService UserService,
	auditLogService AuditLogService,
) *AccountGroupService {
	return &AccountGroupService{
		accountGroupRepository: accountGroup,
		transactor:             transactor,
		userService:            userService,
		auditLogService:        auditLogService,
	}
}
