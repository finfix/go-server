package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/modules/auditLog/model"
	repoModel "server/internal/modules/auditLog/repository/model"
	"server/internal/utils/errors"
)

// GetAuditLogs возвращает записи аудит-лога, доступные пользователю по группам счетов
func (s *AuditLogService) GetAuditLogs(ctx context.Context, req model.GetAuditLogsReq) ([]model.AuditLog, error) {
	ctx, span := tracer.Start(ctx, "GetAuditLogs")
	defer span.End()

	// Получаем группы счетов, доступные пользователю
	accessedAccountGroupIDs, err := s.userToAccountGroupService.GetAccessedAccountGroups(ctx, req.Necessary.UserID)
	if err != nil {
		return nil, err
	}

	var accountGroupIDs []uuid.UUID

	// Если запрошена конкретная группа счетов, проверяем, что она доступна пользователю
	if req.AccountGroupID != nil {
		if !slices.Contains(accessedAccountGroupIDs, *req.AccountGroupID) {
			return nil, errors.Forbidden.New("Access denied").
				WithContextParams(ctx).
				WithParams("UserID", req.Necessary.UserID, "AccountGroupID", *req.AccountGroupID).
				WithCustomHumanText("Вы не имеете доступа к группе счетов с ID = %v", *req.AccountGroupID)
		}
		accountGroupIDs = []uuid.UUID{*req.AccountGroupID}
	} else {
		accountGroupIDs = accessedAccountGroupIDs
	}

	return s.auditLogRepository.GetAuditLogs(ctx, repoModel.GetAuditLogsReq{
		AccountGroupIDs: accountGroupIDs,
		Entity:          req.Entity,
		Method:          req.Method,
		EntityID:        req.EntityID,
		Limit:           req.Limit,
		Offset:          req.Offset,
	})
}
