// Package mcp — инструмент MCP для справочника групп счетов, только чтение
package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/accountGroup/model"
	accountGroupService "server/internal/modules/accountGroup/service"
	mcpTransport "server/internal/utils/mcp"
)

var _ AccountGroupService = new(accountGroupService.AccountGroupService)

// AccountGroupService — интерфейс бизнес-логики групп счетов
type AccountGroupService interface {
	GetAccountGroups(context.Context, model.GetAccountGroupsReq) ([]model.AccountGroup, error)
}

// AccountGroupTools — обработчики инструментов групп счетов
type AccountGroupTools struct {
	accountGroupService AccountGroupService
}

// GetAccountGroupsIn — аргументы get_account_groups
type GetAccountGroupsIn struct{}

// GetAccountGroupsOut — результат get_account_groups
type GetAccountGroupsOut struct {
	AccountGroups []AccountGroupView `json:"accountGroups"` // Группы счетов
}

// AccountGroupView — группа счетов в ответе инструмента
type AccountGroupView struct {
	ID       string `json:"id"`       // Идентификатор
	Name     string `json:"name"`     // Название
	Currency string `json:"currency"` // Валюта группы
	Visible  bool   `json:"visible"`  // Показывается ли в приложении
}

// Register добавляет инструменты групп счетов на MCP-сервер
func Register(server *sdk.Server, accountGroupService AccountGroupService) {
	tools := &AccountGroupTools{accountGroupService: accountGroupService}
	mcpTransport.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "get_account_groups",
		Description: "Группы счетов, доступные пользователю (личные и общие, например семейная), с их валютами",
		Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}, //nolint:exhaustruct
	}, tools.GetAccountGroups)
}
