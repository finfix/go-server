package model

import "github.com/finfix/go-server-grpc/proto"

// GetAuditLogsRes - ответ на получение записей аудит-лога
type GetAuditLogsRes struct {
	AuditLogs []AuditLog
}

// ConvertToProto преобразует GetAuditLogsRes в proto-формат
func (s *GetAuditLogsRes) ConvertToProto() (res *proto.GetAuditLogsResponse, err error) {

	protoAuditLogs := make([]*proto.AuditLog, 0, len(s.AuditLogs))
	for _, auditLog := range s.AuditLogs {
		protoAuditLog, err := auditLog.ConvertToProto()
		if err != nil {
			return res, err
		}
		protoAuditLogs = append(protoAuditLogs, protoAuditLog)
	}

	return &proto.GetAuditLogsResponse{
		Error:     nil,
		AuditLogs: protoAuditLogs,
	}, nil
}
