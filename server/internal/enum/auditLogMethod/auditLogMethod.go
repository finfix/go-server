package auditLogMethod

import (
	"context"

	"pkg/maps"

	"server/internal/utils/errors"

	"github.com/finfix/go-server-grpc/proto"
)

// AuditLogMethod - тип изменяющего действия, зафиксированного в аудит-логе
type AuditLogMethod string

// enum:"create,update,delete"
const (
	Create = AuditLogMethod("create")
	Update = AuditLogMethod("update")
	Delete = AuditLogMethod("delete")
)

// Validate проверяет, что метод аудит-лога принадлежит известному набору значений
func (m AuditLogMethod) Validate(ctx context.Context) error {
	switch m {
	case Create, Update, Delete:
	default:
		return errors.BadRequest.New("Unknown audit log method").
			WithContextParams(ctx).
			WithParams("method", m)
	}
	return nil
}

// mappingProtoToModel содержит соответствие между значениями proto.AuditLogMethod и AuditLogMethod
var mappingProtoToModel = map[proto.AuditLogMethod]AuditLogMethod{
	proto.AuditLogMethod_Create: Create,
	proto.AuditLogMethod_Update: Update,
	proto.AuditLogMethod_Delete: Delete,
}

// ConvertToProto преобразует AuditLogMethod в proto.AuditLogMethod
func (m AuditLogMethod) ConvertToProto() (auditLogMethod proto.AuditLogMethod, err error) {

	// Разворачиваем мапу
	mappingModelToProto, err := maps.Revert(mappingProtoToModel)
	if err != nil {
		return 0, err
	}

	// Получаем значение
	protoAuditLogMethod, ok := mappingModelToProto[m]
	if !ok {
		return protoAuditLogMethod, errors.BadRequest.New("AuditLogMethod undefined")
	}

	return protoAuditLogMethod, nil
}

type ProtoAuditLogMethod struct {
	proto.AuditLogMethod
}

// ConvertToModel преобразует ProtoAuditLogMethod в AuditLogMethod
func (p ProtoAuditLogMethod) ConvertToModel() (auditLogMethod AuditLogMethod, err error) {

	// Проверяем наличие значения
	auditLogMethod, ok := mappingProtoToModel[p.AuditLogMethod]
	if !ok {
		return auditLogMethod, errors.BadRequest.New("AuditLogMethod undefined")
	}

	return auditLogMethod, nil
}
