package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"

	accountGroupService "server/internal/modules/accountGroup/service"
	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	tagModel "server/internal/modules/tag/model"
	tagRepository "server/internal/modules/tag/repository"
	tagRepoModel "server/internal/modules/tag/repository/model"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
	userToAccountGroupService "server/internal/modules/userToAccountGroup/service"
)

var tracer = otel.Tracer("/server/internal/modules/tag/service")

type TagService struct {
	tagRepository             TagRepository
	generalRepository         Transactor
	userToAccountGroupService UserToAccountGroupService
	accountGroupService       AccountGroupService
	auditLogService           AuditLogService
	userService               UserService
}

func NewTagService(
	tagRepository TagRepository,
	generalRepository Transactor,
	userToAccountGroupService UserToAccountGroupService,
	accountGroupService AccountGroupService,
	auditLogService AuditLogService,
	userService UserService,
) *TagService {
	return &TagService{
		tagRepository:             tagRepository,
		generalRepository:         generalRepository,
		userToAccountGroupService: userToAccountGroupService,
		accountGroupService:       accountGroupService,
		auditLogService:           auditLogService,
		userService:               userService,
	}
}

var _ UserToAccountGroupService = new(userToAccountGroupService.UserToAccountGroupService)

type UserToAccountGroupService interface {
	GetAccessedAccountGroups(ctx context.Context, userID uuid.UUID) (accesses []uuid.UUID, err error)
}

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
	WithSyncGate(ctx context.Context, userID uuid.UUID, deviceID string, deviceSyncGate transactor.DeviceSyncGate, auditLogChangeChecker transactor.AuditLogChangeChecker, mutate func(ctxTx context.Context) (auditLogID uint32, err error)) error
}

var _ TagRepository = new(tagRepository.TagRepository)

type TagRepository interface {
	CreateTag(context.Context, tagRepoModel.CreateTagReq) error
	UpdateTag(context.Context, tagModel.UpdateTagReq) error
	DeleteTag(ctx context.Context, id, userID uuid.UUID) error
	GetTags(context.Context, tagModel.GetTagsReq) (res []tagModel.Tag, err error)

	GetTagsToTransactions(ctx context.Context, req tagModel.GetTagsToTransactionsReq) ([]tagModel.TagToTransaction, error)

	CheckAccess(ctx context.Context, accountGroupIDs, tagIDs []uuid.UUID) error
}

var _ AccountGroupService = new(accountGroupService.AccountGroupService)

type AccountGroupService interface {
	CheckAccess(ctx context.Context, userID uuid.UUID, accountGroupIDs []uuid.UUID) error
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
