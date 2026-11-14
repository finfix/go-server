package mcp

import (
	"context"

	"server/internal/utils/necessary"
)

// Necessary возвращает пользователя и устройство из контекста запроса для сервисных моделей
func Necessary(ctx context.Context) (necessary.NecessaryUserInformation, error) {
	var dest struct {
		Necessary necessary.NecessaryUserInformation
	}
	if err := necessary.ParseNecessary(ctx, &dest); err != nil {
		return dest.Necessary, err
	}
	return dest.Necessary, nil
}
