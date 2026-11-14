package service

import (
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

// rateAt возвращает курс валюты на момент: последнюю точку не позже момента. Если таких нет —
// первую точку после, а валюта помечается приблизительной. Без истории — текущий курс, тоже
// приблизительно
func (d dataset) rateAt(currency string, at time.Time) (decimal.Decimal, bool) {
	points := d.rateHistory[currency]
	if len(points) == 0 {
		rate, ok := d.rates[currency]
		if ok {
			d.approximate[currency] = true
		}
		return rate, ok && rate.IsPositive()
	}

	// Первая точка позже момента; предыдущая — действующая
	index := sort.Search(len(points), func(i int) bool { return points[i].at.After(at) })
	if index == 0 {
		d.approximate[currency] = true
		return points[0].rate, true
	}
	return points[index-1].rate, true
}
