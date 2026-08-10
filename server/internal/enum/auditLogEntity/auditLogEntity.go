package auditLogEntity

import (
	"context"

	"server/internal/utils/errors"
)

// AuditLogEntity - тип сущности, изменение которой зафиксировано в аудит-логе
type AuditLogEntity string

// enum:"transaction,account,accountGroup,tag,user,currency"
const (
	Transaction  = AuditLogEntity("transaction")
	Account      = AuditLogEntity("account")
	AccountGroup = AuditLogEntity("accountGroup")
	Tag          = AuditLogEntity("tag")
	User         = AuditLogEntity("user")
	Currency     = AuditLogEntity("currency")
)

// Validate проверяет, что тип сущности аудит-лога принадлежит известному набору значений
func (e AuditLogEntity) Validate(ctx context.Context) error {
	switch e {
	case Transaction, Account, AccountGroup, Tag, User, Currency:
	default:
		return errors.BadRequest.New("Unknown audit log entity").
			WithContextParams(ctx).
			WithParams("entity", e)
	}
	return nil
}
