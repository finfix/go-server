package transactionType

import (
	"context"

	"server/internal/utils/errors"
)

type Type string

// enum:"consumption,income,transfer"
const (
	Transfer    = Type("transfer")
	Consumption = Type("consumption")
	Balancing   = Type("balancing")
	Income      = Type("income")
)

func (t *Type) Validate(ctx context.Context) error {
	if t == nil {
		return nil
	}
	switch *t {
	case Transfer, Consumption, Balancing, Income:
	default:
		return errors.BadRequest.New("Unknown transaction type").
			WithContextParams(ctx).
			SkipThisCall().
			WithParams("type", *t).
			WithCustomHumanText("Неизвестный тип транзакции")
	}
	return nil
}
