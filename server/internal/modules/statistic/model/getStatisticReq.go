package model

import (
	"server/internal/utils/necessary"
)

// GetStatisticReq — запрос рядов графика
type GetStatisticReq struct {
	Necessary            necessary.NecessaryUserInformation
	ChartType            ChartType // Вид графика
	Period               Period    // Размер бакета
	GroupBy              GroupBy   // Разбивка доходов и расходов
	AggregateIntoParents bool      // Сворачивать дочерние счета в родителя с is_parent
	Currency             string    // Валюта ответа, пусто — по умолчанию
	Filters              Filters   // Фильтры транзакций
}
