package logger

import (
	"context"

	"pkg/contextKeys"
	contextKeys2 "server/internal/utils/contextKeys"
)

// GetUserInfo извлекает дополнительную информацию из контекста
func GetUserInfo(ctx context.Context) *UserInfo {

	var userInfo UserInfo

	if ctx == nil {
		return nil
	}

	claims, err := contextKeys2.GetAuthClaims(ctx)
	if err == nil {
		userInfo.UserID = &claims.UserID
		userInfo.DeviceID = &claims.DeviceID
	}

	userInfo.RequestID = contextKeys.GetRequestID(ctx)

	return &userInfo
}
