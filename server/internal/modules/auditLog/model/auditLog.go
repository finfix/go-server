package model

import (
	"time"

	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
)

// AuditLog - запись аудит-лога об изменяющем действии пользователя
type AuditLog struct {
	ID             uint32                        // Числовой идентификатор записи
	Entity         auditLogEntity.AuditLogEntity // Название сущности (transaction, account и т.д.)
	Method         auditLogMethod.AuditLogMethod // Метод изменения (create, update, delete)
	EntityID       string                        // Идентификатор сущности (без FK)
	SnapshotBefore []byte                        // Слепок сущности до изменения в формате JSON
	SnapshotAfter  []byte                        // Слепок сущности после изменения в формате JSON
	UserID         uuid.UUID                     // Идентификатор пользователя, совершившего действие
	DeviceID       string                        // Идентификатор устройства, с которого совершено действие
	AccountGroupID *uuid.UUID                    // Идентификатор группы счетов, в которой изменена сущность (nil для глобальных сущностей)
	DatetimeCreate time.Time                     // Дата и время изменения
}
