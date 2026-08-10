package model

import (
	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	"github.com/finfix/go-server-grpc/proto"
)

// GetAuditLogsReq - запрос на получение записей аудит-лога, доступных пользователю
type GetAuditLogsReq struct {
	Necessary      necessary.NecessaryUserInformation
	AccountGroupID *uuid.UUID                     `json:"accountGroupID"` // Фильтр по группе счетов (должна быть доступна пользователю)
	Entity         *auditLogEntity.AuditLogEntity `json:"entity"`         // Фильтр по названию сущности
	Method         *auditLogMethod.AuditLogMethod `json:"method"`         // Фильтр по методу изменения
	EntityID       *string                        `json:"entityID"`       // Фильтр по идентификатору сущности
	Limit          *uint32                        `json:"limit"`          // Ограничение количества записей
	Offset         *uint32                        `json:"offset"`         // Смещение выборки
}

// ProtoGetAuditLogsReq wrapper for proto request
type ProtoGetAuditLogsReq struct {
	*proto.GetAuditLogsRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoGetAuditLogsReq) ConvertToModel() (res GetAuditLogsReq, err error) {
	if p.GetAuditLogsRequest == nil {
		return res, errors.BadRequest.New("GetAuditLogsRequest is required")
	}

	// Парсим опциональный фильтр по группе счетов
	var accountGroupID *uuid.UUID
	if p.AccountGroupID != nil {
		parsedAccountGroupID, err := uuid.FromBytes(p.AccountGroupID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		accountGroupID = &parsedAccountGroupID
	}

	// Парсим опциональный фильтр по названию сущности
	var entity *auditLogEntity.AuditLogEntity
	if p.Entity != nil {
		parsedEntity, err := auditLogEntity.ProtoAuditLogEntity{AuditLogEntity: *p.Entity}.ConvertToModel()
		if err != nil {
			return res, err
		}
		entity = &parsedEntity
	}

	// Парсим опциональный фильтр по методу изменения
	var method *auditLogMethod.AuditLogMethod
	if p.Method != nil {
		parsedMethod, err := auditLogMethod.ProtoAuditLogMethod{AuditLogMethod: *p.Method}.ConvertToModel()
		if err != nil {
			return res, err
		}
		method = &parsedMethod
	}

	return GetAuditLogsReq{
		AccountGroupID: accountGroupID,
		Entity:         entity,
		Method:         method,
		EntityID:       p.EntityID,
		Limit:          p.Limit,
		Offset:         p.Offset,
	}, nil
}
