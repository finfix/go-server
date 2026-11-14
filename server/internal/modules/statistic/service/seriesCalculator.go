package service

import (
	"sort"
	"time"

	"server/internal/modules/statistic/model"
)

// seriesCalculator раскладывает суммы транзакций по рядам и бакетам
type seriesCalculator struct {
	data                 dataset                  // Данные групп
	buckets              []time.Time              // Начала бакетов
	endOfPeriod          time.Time                // Первый день после периода
	aggregateIntoParents bool                     // Сворачивать дочерние счета в родителя с is_parent
	series               map[string]*model.Series // Ряды по ключу
}

// bucketIndex возвращает индекс бакета даты, -1 — вне периода
func (c seriesCalculator) bucketIndex(date time.Time) int {
	if date.Before(c.buckets[0]) || !date.Before(c.endOfPeriod) {
		return -1
	}
	return sort.Search(len(c.buckets), func(i int) bool { return c.buckets[i].After(date) }) - 1
}
