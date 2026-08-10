package model

import (
	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
)

// TrackMutationReq - запрос на фиксацию изменяющего действия в аудит-логе
type TrackMutationReq struct {
	Entity         auditLogEntity.AuditLogEntity // Название сущности (transaction, account и т.д.)
	Method         auditLogMethod.AuditLogMethod // Метод изменения (create, update, delete)
	EntityID       string                        // Идентификатор сущности (без FK)
	Before         any                           // Слепок сущности до изменения (nil при создании)
	After          any                           // Слепок сущности после изменения (nil при удалении)
	UserID         uuid.UUID                     // Идентификатор пользователя, совершившего действие
	DeviceID       string                        // Идентификатор устройства, с которого совершено действие
	AccountGroupID *uuid.UUID                    // Идентификатор группы счетов, в которой изменена сущность (nil для глобальных сущностей)
}
