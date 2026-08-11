package model

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AccountBudget - версия бюджета счета, действующая с определенной даты
type AccountBudget struct {
	ID              uuid.UUID       `db:"id"`                 // Идентификатор версии бюджета
	AccountID       uuid.UUID       `db:"account_id"`         // Идентификатор счета
	AccountGroupID  uuid.UUID       `db:"account_group_id"`   // Идентификатор группы счетов
	Amount          decimal.Decimal `db:"amount"`             // Сумма бюджета
	FixedSum        decimal.Decimal `db:"fixed_sum"`          // Фиксированная сумма
	DaysOffset      uint32          `db:"days_offset"`        // Смещение в днях
	GradualFilling  bool            `db:"gradual_filling"`    // Заполняется ли бюджет постепенно
	EffectiveFrom   time.Time       `db:"effective_from"`     // Дата, с которой действует эта версия бюджета
	CreatedByUserID uuid.UUID       `db:"created_by_user_id"` // Идентификатор пользователя, создавшего версию
	DatetimeCreate  time.Time       `db:"datetime_create"`    // Дата и время создания версии
}

// ConvertToProto преобразует AccountBudget в proto-формат
func (a AccountBudget) ConvertToProto() *proto.AccountBudget {
	return &proto.AccountBudget{
		Id:              a.ID[:],
		AccountID:       a.AccountID[:],
		AccountGroupID:  a.AccountGroupID[:],
		Amount:          a.Amount.InexactFloat64(),
		FixedSum:        a.FixedSum.InexactFloat64(),
		DaysOffset:      a.DaysOffset,
		GradualFilling:  a.GradualFilling,
		EffectiveFrom:   timestamppb.New(a.EffectiveFrom),
		CreatedByUserID: a.CreatedByUserID[:],
		DatetimeCreate:  timestamppb.New(a.DatetimeCreate),
	}
}
