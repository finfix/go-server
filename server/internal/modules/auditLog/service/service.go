package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"

	"server/internal/modules/auditLog/model"
	"server/internal/modules/auditLog/repository"
	repoModel "server/internal/modules/auditLog/repository/model"
	userToAccountGroupService "server/internal/modules/userToAccountGroup/service"
)

var tracer = otel.Tracer("/server/internal/modules/auditLog/service")

// AuditLogService - сервис аудит-лога изменяющих действий пользователей
type AuditLogService struct {
	auditLogRepository        AuditLogRepository
	userToAccountGroupService UserToAccountGroupService
}

var _ AuditLogRepository = new(repository.AuditLogRepository)

// AuditLogRepository - интерфейс репозитория аудит-лога
type AuditLogRepository interface {
	CreateAuditLog(context.Context, repoModel.CreateAuditLogReq) (uint32, error)
	GetAuditLogs(context.Context, repoModel.GetAuditLogsReq) ([]model.AuditLog, error)
	GetAuditLogsSince(ctx context.Context, accountGroupIDs []uuid.UUID, userID uuid.UUID, sinceID uint32) ([]model.AuditLog, error)
	HasAuditLogsSince(ctx context.Context, accountGroupIDs []uuid.UUID, sinceID uint32) (bool, error)
}

var _ UserToAccountGroupService = new(userToAccountGroupService.UserToAccountGroupService)

// UserToAccountGroupService - интерфейс сервиса связей пользователей с группами счетов
type UserToAccountGroupService interface {
	GetAccessedAccountGroups(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

// NewAuditLogService создает новый сервис аудит-лога
func NewAuditLogService(
	auditLogRepository AuditLogRepository,
	userToAccountGroupService UserToAccountGroupService,
) *AuditLogService {
	return &AuditLogService{
		auditLogRepository:        auditLogRepository,
		userToAccountGroupService: userToAccountGroupService,
	}
}
