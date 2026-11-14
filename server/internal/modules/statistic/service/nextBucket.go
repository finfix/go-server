package service

import (
	"time"

	"server/internal/modules/statistic/model"
)

// nextBucket возвращает начало следующего бакета
func nextBucket(start time.Time, period model.Period) time.Time {
	switch period {
	case model.PeriodWeek:
		return start.AddDate(0, 0, 7)
	case model.PeriodMonth:
		return start.AddDate(0, 1, 0)
	case model.PeriodQuarter:
		return start.AddDate(0, 3, 0)
	case model.PeriodYear:
		return start.AddDate(1, 0, 0)
	}
	return start.AddDate(0, 0, 1)
}
