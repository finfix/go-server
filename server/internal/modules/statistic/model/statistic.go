package model

import "time"

// Statistic — ряды графика с границами бакетов
type Statistic struct {
	Currency string      // Валюта всех сумм
	DateFrom time.Time   // Начало первого бакета
	DateTo   time.Time   // Конец периода включительно
	Buckets  []time.Time // Начала бакетов по возрастанию
	Series   []Series    // Ряды
}
