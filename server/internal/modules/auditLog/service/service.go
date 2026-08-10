package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"server/internal/modules/auditLog/repository"
	repoModel "server/internal/modules/auditLog/repository/model"
)

var tracer = otel.Tracer("/server/internal/modules/auditLog/service")

// AuditLogService - сервис аудит-лога изменяющих действий пользователей
type AuditLogService struct {
	auditLogRepository AuditLogRepository
}

var _ AuditLogRepository = new(repository.AuditLogRepository)

// AuditLogRepository - интерфейс репозитория аудит-лога
type AuditLogRepository interface {
	CreateAuditLog(context.Context, repoModel.CreateAuditLogReq) error
}

// NewAuditLogService создает новый сервис аудит-лога
func NewAuditLogService(
	auditLogRepository AuditLogRepository,
) *AuditLogService {
	return &AuditLogService{
		auditLogRepository: auditLogRepository,
	}
}
