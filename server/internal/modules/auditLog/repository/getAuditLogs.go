package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/modules/auditLog/model"
	"server/internal/modules/auditLog/repository/auditLogDDL"
	repoModel "server/internal/modules/auditLog/repository/model"
)

// GetAuditLogs возвращает записи аудит-лога по фильтрам
func (r *AuditLogRepository) GetAuditLogs(ctx context.Context, req repoModel.GetAuditLogsReq) (auditLogs []model.AuditLog, err error) {
	ctx, span := tracer.Start(ctx, "GetAuditLogs")
	defer span.End()

	// Нет ни одной доступной группы счетов - нет ни одной видимой записи
	if len(req.AccountGroupIDs) == 0 {
		return auditLogs, nil
	}

	filtersEq := sq.Eq{
		auditLogDDL.ColumnAccountGroupID: req.AccountGroupIDs,
	}

	if req.Entity != nil {
		filtersEq[auditLogDDL.ColumnEntity] = *req.Entity
	}
	if req.Method != nil {
		filtersEq[auditLogDDL.ColumnMethod] = *req.Method
	}
	if req.EntityID != nil {
		filtersEq[auditLogDDL.ColumnEntityID] = *req.EntityID
	}

	q := sq.
		Select(ddlHelper.SelectAll).
		From(auditLogDDL.Table).
		Where(filtersEq).
		OrderBy(ddlHelper.Desc(auditLogDDL.ColumnID))

	if req.Limit != nil {
		q = q.Limit(uint64(*req.Limit))
	}
	if req.Offset != nil {
		q = q.Offset(uint64(*req.Offset))
	}

	// Получаем записи аудит-лога
	return auditLogs, r.db.Select(ctx, &auditLogs, q)
}
