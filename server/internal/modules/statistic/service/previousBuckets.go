package service

import (
	"time"

	"server/internal/modules/statistic/model"
)

// previousBuckets отступает от начала бакета на count бакетов назад
func previousBuckets(start time.Time, period model.Period, count int) time.Time {
	switch period {
	case model.PeriodWeek:
		return start.AddDate(0, 0, -7*count)
	case model.PeriodMonth:
		return start.AddDate(0, -count, 0)
	case model.PeriodQuarter:
		return start.AddDate(0, -3*count, 0)
	case model.PeriodYear:
		return start.AddDate(-count, 0, 0)
	}
	return start.AddDate(0, 0, -count)
}
