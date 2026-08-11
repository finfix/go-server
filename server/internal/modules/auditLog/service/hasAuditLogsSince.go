package service

import (
	"context"

	"github.com/google/uuid"
)

// HasAuditLogsSince легковесно проверяет, есть ли у пользователя (по доступным ему группам счетов)
// изменения после sinceID - используется для отсечения мутаций от несинхронизированных устройств
func (s *AuditLogService) HasAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) (bool, error) {
	ctx, span := tracer.Start(ctx, "HasAuditLogsSince")
	defer span.End()

	// Получаем группы счетов, доступные пользователю
	accessedAccountGroupIDs, err := s.userToAccountGroupService.GetAccessedAccountGroups(ctx, userID)
	if err != nil {
		return false, err
	}

	return s.auditLogRepository.HasAuditLogsSince(ctx, accessedAccountGroupIDs, sinceID)
}
