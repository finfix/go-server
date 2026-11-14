package service

import "time"

// endOfDay — последний момент дня даты по UTC: операция без времени берет курс дня, который
// провайдер считает на конец дня
func endOfDay(date time.Time) time.Time {
	return startOf(date, "day").AddDate(0, 0, 1).Add(-time.Nanosecond)
}
