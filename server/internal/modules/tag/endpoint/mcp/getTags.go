package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/tag/model"
	mcpTransport "server/internal/utils/mcp"
)

// GetTags справочник тегов
func (t *TagTools) GetTags(ctx context.Context, _ *sdk.CallToolRequest, in GetTagsIn) (*sdk.CallToolResult, GetTagsOut, error) {
	var out GetTagsOut

	// Собираем запрос из аргументов и контекста
	necessary, err := mcpTransport.Necessary(ctx)
	if err != nil {
		return nil, out, err
	}
	accountGroupIDs, err := mcpTransport.ParseUUIDs("accountGroupIDs", in.AccountGroupIDs)
	if err != nil {
		return nil, out, err
	}

	// Вызываем метод сервиса
	tags, err := t.tagService.GetTags(ctx, model.GetTagsReq{Necessary: necessary, AccountGroupIDs: accountGroupIDs}) //nolint:exhaustruct
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Конвертируем ответ
	out.Tags = make([]TagView, 0, len(tags))
	for _, tag := range tags {
		out.Tags = append(out.Tags, TagView{ID: tag.ID.String(), Name: tag.Name, AccountGroupID: tag.AccountGroupID.String()})
	}
	return nil, out, nil
}
