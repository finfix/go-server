package accountType

import (
	"context"

	"server/internal/utils/errors"
)

type Type string

// enums:"regular,expense,debt,income,balancing"
const (
	Regular   Type = "regular"
	Expense   Type = "expense"
	Debt      Type = "debt"
	Earnings  Type = "earnings"
	Balancing Type = "balancing"
)

func (t *Type) Validate(ctx context.Context) error {
	if t == nil {
		return nil
	}
	switch *t {
	case Earnings, Expense, Debt, Regular, Balancing:
	default:
		return errors.BadRequest.New("Unknown account type").
			SkipThisCall().
			WithParams("type", *t).
			WithCustomHumanText("Неизвестный тип счета")
	}
	return nil
}
