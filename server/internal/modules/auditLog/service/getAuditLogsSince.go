package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/modules/auditLog/model"
)

// GetAuditLogsSince возвращает все записи аудит-лога с id больше sinceID, видимые пользователю: по
// доступным ему группам счетов, а также изменения собственного профиля и изменения валют (глобальные)
func (s *AuditLogService) GetAuditLogsSince(ctx context.Context, userID uuid.UUID, sinceID uint32) ([]model.AuditLog, error) {
	ctx, span := tracer.Start(ctx, "GetAuditLogsSince")
	defer span.End()

	// Получаем группы счетов, доступные пользователю
	accessedAccountGroupIDs, err := s.userToAccountGroupService.GetAccessedAccountGroups(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.auditLogRepository.GetAuditLogsSince(ctx, accessedAccountGroupIDs, userID, sinceID)
}
