package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shopspring/decimal"

	"server/internal/modules/statistic/model"
	mcpTransport "server/internal/utils/mcp"
)

// GetStatisticIn — аргументы get_statistic
type GetStatisticIn struct {
	ChartType            string    `json:"chartType" jsonschema:"income_expense, income, expense, balance_detailed, balance_total или delta"`
	Period               string    `json:"period" jsonschema:"Бакет: day, week, month, quarter или year"`
	GroupBy              string    `json:"groupBy,omitempty" jsonschema:"Разбивка income/expense: account (по умолчанию) или tag"`
	AggregateIntoParents *bool     `json:"aggregateIntoParents,omitempty" jsonschema:"Сворачивать дочерние счета в родителя с isParent, по умолчанию true"`
	Currency             string    `json:"currency,omitempty" jsonschema:"Валюта ответа; по умолчанию валюта единственной выбранной группы, иначе валюта пользователя"`
	Filters              FiltersIn `json:"filters,omitzero" jsonschema:"Фильтры транзакций"`
}

// GetStatisticOut — результат get_statistic
type GetStatisticOut struct {
	Currency string       `json:"currency"` // Валюта всех сумм
	DateFrom string       `json:"dateFrom"` // Начало первого бакета
	DateTo   string       `json:"dateTo"`   // Конец периода включительно
	Buckets  []string     `json:"buckets"`  // Начала бакетов, yyyy-MM-dd
	Series   []SeriesView `json:"series"`   // Ряды
}

// SeriesView — ряд графика
type SeriesView struct {
	Key            string            `json:"key"`                     // income, expense, delta, balance, untagged или идентификатор счета/тега
	Name           string            `json:"name"`                    // Название
	AccountID      *string           `json:"accountID,omitempty"`     // Счет ряда
	TagID          *string           `json:"tagID,omitempty"`         // Тег ряда
	Values         []decimal.Decimal `json:"values"`                  // Значения по бакетам
	Sum            decimal.Decimal   `json:"sum"`                     // Сумма по бакетам
	Average        decimal.Decimal   `json:"average"`                 // Среднее по бакетам
	AverageNonZero decimal.Decimal   `json:"averageNonZero"`          // Среднее по ненулевым бакетам
	Min            decimal.Decimal   `json:"min"`                     // Минимум
	Max            decimal.Decimal   `json:"max"`                     // Максимум
	MonthlyBudget  *decimal.Decimal  `json:"monthlyBudget,omitempty"` // Месячный бюджет расходного счета
}

// GetStatistic ряды графика
func (t *StatisticTools) GetStatistic(ctx context.Context, _ *sdk.CallToolRequest, in GetStatisticIn) (*sdk.CallToolResult, GetStatisticOut, error) {
	var out GetStatisticOut

	// Собираем запрос из аргументов и контекста
	necessary, err := mcpTransport.Necessary(ctx)
	if err != nil {
		return nil, out, err
	}
	filters, err := in.Filters.ConvertToModel()
	if err != nil {
		return nil, out, err
	}
	aggregateIntoParents := in.AggregateIntoParents == nil || *in.AggregateIntoParents

	// Вызываем метод сервиса
	statistic, err := t.statisticService.GetStatistic(ctx, model.GetStatisticReq{
		Necessary:            necessary,
		ChartType:            model.ChartType(in.ChartType),
		Period:               model.Period(in.Period),
		GroupBy:              model.GroupBy(in.GroupBy),
		AggregateIntoParents: aggregateIntoParents,
		Currency:             in.Currency,
		Filters:              filters,
	})
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Конвертируем ответ
	out = GetStatisticOut{
		Currency: statistic.Currency,
		DateFrom: statistic.DateFrom.Format(mcpTransport.DateLayout),
		DateTo:   statistic.DateTo.Format(mcpTransport.DateLayout),
		Buckets:  make([]string, 0, len(statistic.Buckets)),
		Series:   make([]SeriesView, 0, len(statistic.Series)),
	}
	for _, bucket := range statistic.Buckets {
		out.Buckets = append(out.Buckets, bucket.Format(mcpTransport.DateLayout))
	}
	for _, series := range statistic.Series {
		out.Series = append(out.Series, SeriesView{
			Key:            series.Key,
			Name:           series.Name,
			AccountID:      mcpTransport.UUIDString(series.AccountID),
			TagID:          mcpTransport.UUIDString(series.TagID),
			Values:         series.Values,
			Sum:            series.Summary.Sum,
			Average:        series.Summary.Average,
			AverageNonZero: series.Summary.AverageNonZero,
			Min:            series.Summary.Min,
			Max:            series.Summary.Max,
			MonthlyBudget:  series.MonthlyBudget,
		})
	}
	return nil, out, nil
}
