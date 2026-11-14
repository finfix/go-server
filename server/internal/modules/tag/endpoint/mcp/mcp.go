// Package mcp — инструмент MCP для справочника тегов, только чтение
package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/modules/tag/model"
	tagService "server/internal/modules/tag/service"
	mcpTransport "server/internal/utils/mcp"
)

var _ TagService = new(tagService.TagService)

// TagService — интерфейс бизнес-логики тегов
type TagService interface {
	GetTags(context.Context, model.GetTagsReq) ([]model.Tag, error)
}

// TagTools — обработчики инструментов тегов
type TagTools struct {
	tagService TagService
}

// GetTagsIn — аргументы get_tags
type GetTagsIn struct {
	AccountGroupIDs []string `json:"accountGroupIDs,omitempty" jsonschema:"Группы счетов, пусто — все доступные"`
}

// GetTagsOut — результат get_tags
type GetTagsOut struct {
	Tags []TagView `json:"tags"` // Теги
}

// TagView — тег в ответе инструмента
type TagView struct {
	ID             string `json:"id"`             // Идентификатор
	Name           string `json:"name"`           // Название
	AccountGroupID string `json:"accountGroupID"` // Группа счетов
}

// Register добавляет инструменты тегов на MCP-сервер
func Register(server *sdk.Server, tagService TagService) {
	tools := &TagTools{tagService: tagService}
	mcpTransport.AddTool(server, &sdk.Tool{ //nolint:exhaustruct
		Name:        "get_tags",
		Description: "Теги (подкатегории) доступных групп счетов",
		Annotations: &sdk.ToolAnnotations{ReadOnlyHint: true}, //nolint:exhaustruct
	}, tools.GetTags)
}
