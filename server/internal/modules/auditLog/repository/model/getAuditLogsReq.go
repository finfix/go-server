package model

import (
	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
)

// GetAuditLogsReq - репозиторная модель запроса на получение записей аудит-лога с уже разрешенными правами доступа
type GetAuditLogsReq struct {
	AccountGroupIDs []uuid.UUID                    // Группы счетов, к которым разрешен доступ
	Entity          *auditLogEntity.AuditLogEntity // Фильтр по названию сущности
	Method          *auditLogMethod.AuditLogMethod // Фильтр по методу изменения
	EntityID        *string                        // Фильтр по идентификатору сущности
	Limit           *uint32                        // Ограничение количества записей
	Offset          *uint32                        // Смещение выборки
}
