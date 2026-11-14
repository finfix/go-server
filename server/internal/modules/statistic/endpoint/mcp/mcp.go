// Package mcp — инструменты MCP для статистики: графики и лента транзакций, только чтение
package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/statistic/model"
	statisticService "server/internal/modules/statistic/service"
	mcpTransport "server/internal/utils/mcp"
)

var _ StatisticService = new(statisticService.StatisticService)

// StatisticService — интерфейс расчета статистики
type StatisticService interface {
	GetStatistic(context.Context, model.GetStatisticReq) (model.Statistic, error)
	GetTransactions(context.Context, model.GetTransactionsReq) (model.TransactionsPage, error)
}

// StatisticTools — обработчики инструментов статистики
type StatisticTools struct {
	statisticService StatisticService
}

// Register добавляет инструменты статистики на MCP-сервер
func Register(server *sdk.Server, statisticService StatisticService) {
	tools := &StatisticTools{statisticService: statisticService}
	readOnly := &sdk.ToolAnnotations{ReadOnlyHint: true} //nolint:exhaustruct

	mcpTransport.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name: "get_statistic",
		Description: "Ряды графика по бакетам периода, как экран статистики приложения. chartType: income_expense — доходы и расходы двумя рядами; " +
			"income / expense — с разбивкой groupBy (account — по счетам-категориям, tag — по тегам; транзакция с несколькими тегами попадает в каждый); " +
			"balance_detailed — остаток каждого счета regular/debt на конец бакета (от первой транзакции); balance_total — сумма остатков; delta — доходы минус расходы. " +
			"Без dateFrom — последние 30 дней / 8 недель / 6 месяцев / 4 квартала / 3 года до dateTo (по умолчанию сегодня). " +
			"К каждому ряду — сводка sum/average/averageNonZero/min/max по бакетам, к расходным счетам — действующий месячный бюджет. " +
			"Фильтр types к графику не применяется; для balance_* учитываются только accountIDs, excludedAccountIDs, currencies и accountGroupIDs",
		Annotations: readOnly,
	}, tools.GetStatistic)
	mcpTransport.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name: "get_transactions",
		Description: "Лента транзакций, как вкладка «Транзакции»: от новых к старым, по дням с суммой расходов за день по выборке. " +
			"Итоги totalCount / totalConsumption / totalIncome — по всей выборке, не по странице. Страница — limit (по умолчанию 100, максимум 500) и offset; nextOffset — следующая страница",
		Annotations: readOnly,
	}, tools.GetTransactions)
}
