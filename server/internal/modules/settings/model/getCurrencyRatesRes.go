package model

// GetCurrencyRatesRes — страница истории курсов
type GetCurrencyRatesRes struct {
	Rates   []CurrencyRate // Точки по возрастанию момента
	HasMore bool           // Есть ещё точки после последней
}
