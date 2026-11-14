package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shopspring/decimal"

	"server/internal/modules/statistic/model"
	mcpTransport "server/internal/utils/mcp"
)

// GetTransactionsIn — аргументы get_transactions
type GetTransactionsIn struct {
	Currency string    `json:"currency,omitempty" jsonschema:"Валюта итогов; по умолчанию валюта единственной выбранной группы, иначе валюта пользователя"`
	Filters  FiltersIn `json:"filters,omitzero" jsonschema:"Фильтры транзакций"`
	Offset   int       `json:"offset,omitempty" jsonschema:"Сколько транзакций пропустить: nextOffset прошлой страницы"`
	Limit    int       `json:"limit,omitempty" jsonschema:"Размер страницы, по умолчанию 100, максимум 500"`
}

// GetTransactionsOut — результат get_transactions
type GetTransactionsOut struct {
	Currency              string          `json:"currency"`                        // Валюта итогов
	TotalCount            int             `json:"totalCount"`                      // Транзакций в выборке
	TotalConsumption      decimal.Decimal `json:"totalConsumption"`                // Расходы по выборке
	TotalIncome           decimal.Decimal `json:"totalIncome"`                     // Доходы по выборке
	NextOffset            *int            `json:"nextOffset,omitempty"`            // Смещение следующей страницы
	ApproximateCurrencies []string        `json:"approximateCurrencies,omitempty"` // Валюты без курса до даты операции: взят ближайший после
	Days                  []DayView       `json:"days"`                            // Транзакции страницы по дням
}

// DayView — день ленты
type DayView struct {
	Date             string            `json:"date"`             // День, yyyy-MM-dd
	ConsumptionTotal decimal.Decimal   `json:"consumptionTotal"` // Расходы за день по выборке
	Transactions     []TransactionView `json:"transactions"`     // Транзакции дня на странице
}

// TransactionView — транзакция ленты
type TransactionView struct {
	ID              string          `json:"id"`               // Идентификатор
	Type            string          `json:"type"`             // consumption, income, transfer или balancing
	AccountFromID   string          `json:"accountFromID"`    // Счет списания
	AccountFromName string          `json:"accountFromName"`  // Название счета списания
	AmountFrom      decimal.Decimal `json:"amountFrom"`       // Сумма списания в валюте счета
	CurrencyFrom    string          `json:"currencyFrom"`     // Валюта счета списания
	AccountToID     string          `json:"accountToID"`      // Счет пополнения
	AccountToName   string          `json:"accountToName"`    // Название счета пополнения
	AmountTo        decimal.Decimal `json:"amountTo"`         // Сумма пополнения в валюте счета
	CurrencyTo      string          `json:"currencyTo"`       // Валюта счета пополнения
	Note            string          `json:"note,omitempty"`   // Заметка
	TagIDs          []string        `json:"tagIDs,omitempty"` // Теги
	AccountGroupID  string          `json:"accountGroupID"`   // Группа счетов
	CreatedByUserID string          `json:"createdByUserID"`  // Автор
}

// GetTransactions лента транзакций
func (t *StatisticTools) GetTransactions(ctx context.Context, _ *sdk.CallToolRequest, in GetTransactionsIn) (*sdk.CallToolResult, GetTransactionsOut, error) {
	var out GetTransactionsOut

	// Собираем запрос из аргументов и контекста
	necessary, err := mcpTransport.Necessary(ctx)
	if err != nil {
		return nil, out, err
	}
	filters, err := in.Filters.ConvertToModel()
	if err != nil {
		return nil, out, err
	}

	// Вызываем метод сервиса
	page, err := t.statisticService.GetTransactions(ctx, model.GetTransactionsReq{
		Necessary: necessary,
		Currency:  in.Currency,
		Filters:   filters,
		Offset:    in.Offset,
		Limit:     in.Limit,
	})
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Конвертируем ответ
	out = GetTransactionsOut{
		Currency:              page.Currency,
		TotalCount:            page.TotalCount,
		TotalConsumption:      page.TotalConsumption,
		TotalIncome:           page.TotalIncome,
		NextOffset:            page.NextOffset,
		ApproximateCurrencies: page.ApproximateCurrencies,
		Days:                  make([]DayView, 0, len(page.Days)),
	}
	for _, day := range page.Days {
		view := DayView{Date: day.Date.Format(mcpTransport.DateLayout), ConsumptionTotal: day.ConsumptionTotal, Transactions: make([]TransactionView, 0, len(day.Transactions))}
		for _, transaction := range day.Transactions {
			view.Transactions = append(view.Transactions, NewTransactionView(transaction))
		}
		out.Days = append(out.Days, view)
	}
	return nil, out, nil
}
