package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	pendingLinkedTransferModel "server/internal/modules/pendingLinkedTransfer/model"
	pendingLinkedTransferRepository "server/internal/modules/pendingLinkedTransfer/repository"
	repoModel "server/internal/modules/pendingLinkedTransfer/repository/model"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
)

var tracer = otel.Tracer("/server/internal/modules/pendingLinkedTransfer/service")

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithSyncGate(ctx context.Context, userID uuid.UUID, deviceID string, deviceSyncGate transactor.DeviceSyncGate, auditLogChangeChecker transactor.AuditLogChangeChecker, mutate func(ctxTx context.Context) (auditLogID uint32, err error)) error
}

var _ PendingLinkedTransferRepository = new(pendingLinkedTransferRepository.PendingLinkedTransferRepository)

type PendingLinkedTransferRepository interface {
	CreatePendingLinkedTransfer(context.Context, repoModel.CreatePendingLinkedTransferReq) error
	GetPendingLinkedTransfers(context.Context, repoModel.GetPendingLinkedTransfersReq) ([]pendingLinkedTransferModel.PendingLinkedTransfer, error)
	UpdatePendingLinkedTransfer(ctx context.Context, id uuid.UUID, req repoModel.UpdatePendingLinkedTransferReq) error
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

// PendingLinkedTransferService - сервис требований довнесения транзакции через счёт-мост.
// Тупая точка правды с CRUD, без проверок доступа/членства в группе — вся оркестрация на
// фронтенде.
type PendingLinkedTransferService struct {
	pendingLinkedTransferRepository PendingLinkedTransferRepository
	transactor                      Transactor
	auditLogService                 AuditLogService
	userService                     UserService
}

func NewPendingLinkedTransferService(
	pendingLinkedTransferRepository PendingLinkedTransferRepository,
	transactor Transactor,
	auditLogService AuditLogService,
	userService UserService,
) *PendingLinkedTransferService {
	return &PendingLinkedTransferService{
		pendingLinkedTransferRepository: pendingLinkedTransferRepository,
		transactor:                      transactor,
		auditLogService:                 auditLogService,
		userService:                     userService,
	}
}
