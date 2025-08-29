package model

import "github.com/shopspring/decimal"

type Currency struct {
	Signatura string          `json:"isoCode" db:"signatura"` // Сигнатура валюты
	Name      string          `json:"name" db:"name"`         // Название валюты
	Symbol    string          `json:"symbol" db:"symbol"`     // Символ валюты
	Rate      decimal.Decimal `json:"rate" db:"rate"`         // Курс валюты
}
