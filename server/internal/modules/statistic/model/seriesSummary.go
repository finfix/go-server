package model

import "github.com/shopspring/decimal"

// SeriesSummary — агрегаты ряда по бакетам
type SeriesSummary struct {
	Sum            decimal.Decimal // Сумма
	Average        decimal.Decimal // Среднее по всем бакетам
	AverageNonZero decimal.Decimal // Среднее по ненулевым бакетам
	Min            decimal.Decimal // Минимум
	Max            decimal.Decimal // Максимум
}
