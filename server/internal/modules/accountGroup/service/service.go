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
	userToAccountGroupService "server/internal/modules/userToAccountGroup/service"

	"github.com/google/uuid"
)

var tracer = otel.Tracer("/server/internal/modules/accountGroup/service")

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
	WithSyncGate(ctx context.Context, userID uuid.UUID, deviceID string, deviceSyncGate transactor.DeviceSyncGate, auditLogChangeChecker transactor.AuditLogChangeChecker, mutate func(ctxTx context.Context) (auditLogID uint32, err error)) error
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

var _ UserToAccountGroupService = new(userToAccountGroupService.UserToAccountGroupService)

type UserToAccountGroupService interface {
	GetAccessedAccountGroups(ctx context.Context, userID uuid.UUID) (accesses []uuid.UUID, err error)
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

type AccountGroupService struct {
	userToAccountGroupService UserToAccountGroupService
	accountGroupRepository    AccountGroupRepository
	transactor                Transactor
	auditLogService           AuditLogService
	userService               UserService
}

func NewAccountGroupService(
	accountGroup AccountGroupRepository,
	transactor Transactor,
	userToAccountGroupService UserToAccountGroupService,
	auditLogService AuditLogService,
	userService UserService,
) *AccountGroupService {
	return &AccountGroupService{
		accountGroupRepository:    accountGroup,
		transactor:                transactor,
		userToAccountGroupService: userToAccountGroupService,
		auditLogService:           auditLogService,
		userService:               userService,
	}
}
