package grpc

import (
	"context"

	"server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"

	"github.com/finfix/go-server-grpc/proto"
)

var _ AuditLogService = new(auditLogService.AuditLogService)

type AuditLogService interface {
	GetAuditLogs(context.Context, model.GetAuditLogsReq) ([]model.AuditLog, error)
}

var _ proto.AuditLogEndpointServer = new(AuditLogEndpoint)

type AuditLogEndpoint struct {
	proto.UnsafeAuditLogEndpointServer
	auditLogService AuditLogService
}

func NewAuditLogEndpoint(auditLogService AuditLogService) *AuditLogEndpoint {
	return &AuditLogEndpoint{
		auditLogService: auditLogService,
	}
}
