package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/finfix/go-server-grpc/proto"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
)

// AuditLog - запись аудит-лога об изменяющем действии пользователя
type AuditLog struct {
	ID             uint32                        `db:"id"`               // Числовой идентификатор записи
	Entity         auditLogEntity.AuditLogEntity `db:"entity"`           // Название сущности (transaction, account и т.д.)
	Method         auditLogMethod.AuditLogMethod `db:"method"`           // Метод изменения (create, update, delete)
	EntityID       string                        `db:"entity_id"`        // Идентификатор сущности (без FK)
	SnapshotBefore []byte                        `db:"snapshot_before"`  // Слепок сущности до изменения в формате JSON
	SnapshotAfter  []byte                        `db:"snapshot_after"`   // Слепок сущности после изменения в формате JSON
	UserID         uuid.UUID                     `db:"user_id"`          // Идентификатор пользователя, совершившего действие
	DeviceID       string                        `db:"device_id"`        // Идентификатор устройства, с которого совершено действие
	AccountGroupID *uuid.UUID                    `db:"account_group_id"` // Идентификатор группы счетов, в которой изменена сущность (nil для глобальных сущностей)
	DatetimeCreate time.Time                     `db:"datetime_create"`  // Дата и время изменения
}

// ConvertToProto преобразует AuditLog в proto-формат
func (a AuditLog) ConvertToProto() (*proto.AuditLog, error) {
	protoEntity, err := a.Entity.ConvertToProto()
	if err != nil {
		return nil, err
	}

	protoMethod, err := a.Method.ConvertToProto()
	if err != nil {
		return nil, err
	}

	var accountGroupID []byte
	if a.AccountGroupID != nil {
		accountGroupID = a.AccountGroupID[:]
	}

	return &proto.AuditLog{
		Id:             a.ID,
		Entity:         protoEntity,
		Method:         protoMethod,
		EntityID:       a.EntityID,
		SnapshotBefore: a.SnapshotBefore,
		SnapshotAfter:  a.SnapshotAfter,
		UserID:         a.UserID[:],
		DeviceID:       a.DeviceID,
		AccountGroupID: accountGroupID,
		DatetimeCreate: timestamppb.New(a.DatetimeCreate),
	}, nil
}
