package model

import (
	"context"

	"server/internal/utils/errors"
)

// Period — размер бакета графика
type Period string

const (
	PeriodDay     Period = "day"     // День
	PeriodWeek    Period = "week"    // Неделя с понедельника
	PeriodMonth   Period = "month"   // Месяц
	PeriodQuarter Period = "quarter" // Квартал
	PeriodYear    Period = "year"    // Год
)

// Validate проверяет, что период известен
func (p Period) Validate(ctx context.Context) error {
	switch p {
	case PeriodDay, PeriodWeek, PeriodMonth, PeriodQuarter, PeriodYear:
		return nil
	}
	return errors.BadRequest.New("unknown period").WithParams("period", string(p)).WithContextParams(ctx)
}
