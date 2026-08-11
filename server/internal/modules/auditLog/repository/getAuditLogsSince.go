package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"pkg/ddlHelper"
	"server/internal/enum/auditLogEntity"
	"server/internal/modules/auditLog/model"
	"server/internal/modules/auditLog/repository/auditLogDDL"
)

// GetAuditLogsSince возвращает все записи аудит-лога с id больше sinceID, видимые пользователю:
// по указанным группам счетов, а также глобальные для этого пользователя записи, не привязанные
// к группам счетов - изменения собственного профиля и изменения валют (общие для всех), отсортированные
// от старых к новым
func (r *AuditLogRepository) GetAuditLogsSince(ctx context.Context, accountGroupIDs []uuid.UUID, userID uuid.UUID, sinceID uint32) (auditLogs []model.AuditLog, err error) {
	ctx, span := tracer.Start(ctx, "GetAuditLogsSince")
	defer span.End()

	visibility := sq.Or{
		sq.And{
			sq.Eq{auditLogDDL.ColumnEntity: string(auditLogEntity.User)},
			sq.Eq{auditLogDDL.ColumnEntityID: userID.String()},
		},
		sq.Eq{auditLogDDL.ColumnEntity: string(auditLogEntity.Currency)},
	}

	if len(accountGroupIDs) != 0 {
		visibility = append(visibility, sq.Eq{auditLogDDL.ColumnAccountGroupID: accountGroupIDs})
	}

	q := sq.
		Select(ddlHelper.SelectAll).
		From(auditLogDDL.Table).
		Where(visibility).
		Where(sq.Gt{
			auditLogDDL.ColumnID: sinceID,
		}).
		OrderBy(auditLogDDL.ColumnID)

	return auditLogs, r.db.Select(ctx, &auditLogs, q)
}
