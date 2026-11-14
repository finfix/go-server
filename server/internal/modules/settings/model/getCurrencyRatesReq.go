package model

import "time"

// GetCurrencyRatesReq — запрос истории курсов
type GetCurrencyRatesReq struct {
	Currencies []string   // Валюты, пусто — все
	DateFrom   *time.Time // Начало периода включительно; nil — с первой точки
	DateTo     *time.Time // Конец периода включительно; nil — до последней точки
}
