package auditLogEntity

import (
	"context"

	"pkg/maps"

	"server/internal/utils/errors"

	"github.com/finfix/go-server-grpc/proto"
)

// AuditLogEntity - тип сущности, изменение которой зафиксировано в аудит-логе
type AuditLogEntity string

// enum:"transaction,account,accountGroup,tag,user,currency,accountBudget"
const (
	Transaction   = AuditLogEntity("transaction")
	Account       = AuditLogEntity("account")
	AccountGroup  = AuditLogEntity("accountGroup")
	Tag           = AuditLogEntity("tag")
	User          = AuditLogEntity("user")
	Currency      = AuditLogEntity("currency")
	AccountBudget = AuditLogEntity("accountBudget")
)

// Validate проверяет, что тип сущности аудит-лога принадлежит известному набору значений
func (e AuditLogEntity) Validate(ctx context.Context) error {
	switch e {
	case Transaction, Account, AccountGroup, Tag, User, Currency, AccountBudget:
	default:
		return errors.BadRequest.New("Unknown audit log entity").
			WithContextParams(ctx).
			WithParams("entity", e)
	}
	return nil
}

// mappingProtoToModel содержит соответствие между значениями proto.AuditLogEntity и AuditLogEntity
var mappingProtoToModel = map[proto.AuditLogEntity]AuditLogEntity{
	proto.AuditLogEntity_Transaction:   Transaction,
	proto.AuditLogEntity_Account:       Account,
	proto.AuditLogEntity_AccountGroup:  AccountGroup,
	proto.AuditLogEntity_Tag:           Tag,
	proto.AuditLogEntity_User:          User,
	proto.AuditLogEntity_Currency:      Currency,
	proto.AuditLogEntity_AccountBudget: AccountBudget,
}

// ConvertToProto преобразует AuditLogEntity в proto.AuditLogEntity
func (e AuditLogEntity) ConvertToProto() (auditLogEntity proto.AuditLogEntity, err error) {

	// Разворачиваем мапу
	mappingModelToProto, err := maps.Revert(mappingProtoToModel)
	if err != nil {
		return 0, err
	}

	// Получаем значение
	protoAuditLogEntity, ok := mappingModelToProto[e]
	if !ok {
		return protoAuditLogEntity, errors.BadRequest.New("AuditLogEntity undefined")
	}

	return protoAuditLogEntity, nil
}

type ProtoAuditLogEntity struct {
	proto.AuditLogEntity
}

// ConvertToModel преобразует ProtoAuditLogEntity в AuditLogEntity
func (p ProtoAuditLogEntity) ConvertToModel() (auditLogEntity AuditLogEntity, err error) {

	// Проверяем наличие значения
	auditLogEntity, ok := mappingProtoToModel[p.AuditLogEntity]
	if !ok {
		return auditLogEntity, errors.BadRequest.New("AuditLogEntity undefined")
	}

	return auditLogEntity, nil
}
