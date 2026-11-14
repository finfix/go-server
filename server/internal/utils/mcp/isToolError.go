package mcp

import (
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// isToolError — результат вызова инструмента с флагом ошибки
func isToolError(result sdk.Result) bool {
	callResult, ok := result.(*sdk.CallToolResult)
	return ok && callResult.IsError
}
