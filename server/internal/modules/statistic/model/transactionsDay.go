package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// TransactionsDay — день ленты: транзакции и сумма расходов за день по выборке
type TransactionsDay struct {
	Date             time.Time         // День
	ConsumptionTotal decimal.Decimal   // Расходы за день
	Transactions     []TransactionItem // Транзакции дня на странице
}
