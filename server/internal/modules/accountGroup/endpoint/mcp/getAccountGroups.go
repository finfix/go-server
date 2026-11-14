package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/accountGroup/model"
	mcpTransport "server/internal/utils/mcp"
)

// GetAccountGroups справочник групп счетов
func (t *AccountGroupTools) GetAccountGroups(ctx context.Context, _ *sdk.CallToolRequest, _ GetAccountGroupsIn) (*sdk.CallToolResult, GetAccountGroupsOut, error) {
	var out GetAccountGroupsOut

	// Собираем запрос из контекста: репозиторий ограничивает группы пользователем
	necessary, err := mcpTransport.Necessary(ctx)
	if err != nil {
		return nil, out, err
	}

	// Вызываем метод сервиса
	groups, err := t.accountGroupService.GetAccountGroups(ctx, model.GetAccountGroupsReq{Necessary: necessary}) //nolint:exhaustruct
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Конвертируем ответ
	out.AccountGroups = make([]AccountGroupView, 0, len(groups))
	for _, group := range groups {
		out.AccountGroups = append(out.AccountGroups, AccountGroupView{ID: group.ID.String(), Name: group.Name, Currency: group.Currency, Visible: group.Visible})
	}
	return nil, out, nil
}
