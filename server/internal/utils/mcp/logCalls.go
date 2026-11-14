package mcp

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	pkgErrors "pkg/errors"
	"pkg/log"
)

// logCalls пишет каждый вызов инструмента: имя, длительность, ошибка. Аргументы не пишем — в них
// финансовые данные пользователя
func logCalls(next sdk.MethodHandler) sdk.MethodHandler {
	return func(ctx context.Context, method string, req sdk.Request) (sdk.Result, error) {
		started := time.Now()
		result, err := next(ctx, method, req)

		name := method
		if params, ok := req.GetParams().(*sdk.CallToolParamsRaw); ok {
			name = fmt.Sprintf("%s %s", method, params.Name)
		}
		duration := time.Since(started).Round(time.Millisecond)

		switch {
		case err != nil:
			log.Warning(fmt.Sprintf("MCP %s за %s: %s", name, duration, pkgErrors.CastError(err).Error()))
		case isToolError(result):
			log.Warning(fmt.Sprintf("MCP %s за %s: инструмент вернул ошибку", name, duration))
		default:
			log.Info(fmt.Sprintf("MCP %s за %s", name, duration))
		}
		return result, err
	}
}
