// Package mcp — инструмент MCP для справочника счетов, только чтение
package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/account/model"
	accountService "server/internal/modules/account/service"
	mcpTransport "server/internal/utils/mcp"
)

var _ AccountService = new(accountService.AccountService)

// AccountService — интерфейс бизнес-логики счетов
type AccountService interface {
	GetAccounts(context.Context, model.GetAccountsReq) ([]model.Account, error)
}

// AccountTools — обработчики инструментов счетов
type AccountTools struct {
	accountService AccountService
}

// Register добавляет инструменты счетов на MCP-сервер
func Register(server *sdk.Server, accountService AccountService) {
	tools := &AccountTools{accountService: accountService}
	mcpTransport.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name: "get_accounts",
		Description: "Счета доступных групп. Типы: regular — деньги, debt — долги, expense — категории расходов, earnings — источники доходов, balancing — корректировки. " +
			"remainder: у regular/debt — текущий остаток, у expense/earnings/balancing — сумма за период dateFrom–dateTo (по умолчанию текущий месяц). " +
			"Иерархия — parentAccountID; isParent — родитель, в которого сворачиваются дочерние",
		Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}, //nolint:exhaustruct
	}, tools.GetAccounts)
}
