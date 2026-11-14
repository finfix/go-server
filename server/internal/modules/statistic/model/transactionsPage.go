package model

import "github.com/shopspring/decimal"

// TransactionsPage — страница ленты транзакций с итогами по всей выборке
type TransactionsPage struct {
	Currency         string            // Валюта итогов
	Days             []TransactionsDay // Транзакции страницы по дням
	TotalCount       int               // Транзакций в выборке
	TotalConsumption decimal.Decimal   // Расходы по выборке
	TotalIncome      decimal.Decimal   // Доходы по выборке
	NextOffset       *int              // Смещение следующей страницы, nil — страница последняя
}
