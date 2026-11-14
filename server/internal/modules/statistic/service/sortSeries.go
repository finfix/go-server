package service

import (
	"slices"

	"server/internal/modules/statistic/model"
)

// sortSeries упорядочивает ряды разбивки: по убыванию суммы, затем по названию
func sortSeries(series []model.Series) {
	slices.SortStableFunc(series, func(a, b model.Series) int {
		if c := b.Summary.Sum.Cmp(a.Summary.Sum); c != 0 {
			return c
		}
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})
}
