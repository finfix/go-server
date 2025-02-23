package accountType

import (
	"context"

	"pkg/errors"
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
		return errors.BadRequest.New(ctx, "Unknown account type",
			errors.SkipThisCallOption(),
			errors.ParamsOption("type", *t),
			errors.HumanTextOption("Неизвестный тип счета"),
		)
	}
	return nil
}
