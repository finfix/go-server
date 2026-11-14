package service

import (
	"time"

	"server/internal/modules/statistic/model"
)

// defaultBucketCount — сколько бакетов показывать без date_from, как в приложении
var defaultBucketCount = map[model.Period]int{
	model.PeriodDay:     30,
	model.PeriodWeek:    8,
	model.PeriodMonth:   6,
	model.PeriodQuarter: 4,
	model.PeriodYear:    3,
}

// maxBuckets — предел бакетов в одном ответе
const maxBuckets = 1000

// startOf возвращает начало бакета, в который попадает дата. Границы по UTC, неделя с понедельника
func startOf(date time.Time, period model.Period) time.Time {
	date = date.UTC()
	day := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	switch period {
	case model.PeriodWeek:
		return day.AddDate(0, 0, -((int(day.Weekday()) + 6) % 7))
	case model.PeriodMonth:
		return time.Date(day.Year(), day.Month(), 1, 0, 0, 0, 0, time.UTC)
	case model.PeriodQuarter:
		return time.Date(day.Year(), day.Month()-(day.Month()-1)%3, 1, 0, 0, 0, 0, time.UTC)
	case model.PeriodYear:
		return time.Date(day.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return day
}
