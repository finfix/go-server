package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/auditLog/model"
	"server/internal/utils/necessary"

	proto "github.com/finfix/go-server-grpc/proto"
)

// GetAuditLogs получение записей аудит-лога, доступных пользователю по группам счетов
func (e *AuditLogEndpoint) GetAuditLogs(ctx context.Context, r *proto.GetAuditLogsRequest) (*proto.GetAuditLogsResponse, error) {
	res := new(proto.GetAuditLogsResponse)

	// Convert proto request to internal model
	req, err := model.ProtoGetAuditLogsReq{GetAuditLogsRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Parse necessary information from context
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return res, err
	}

	// Call service method
	auditLogs, err := e.auditLogService.GetAuditLogs(ctx, req)
	if err != nil {
		return res, err
	}

	_res := model.GetAuditLogsRes{AuditLogs: auditLogs}

	return _res.ConvertToProto()
}
