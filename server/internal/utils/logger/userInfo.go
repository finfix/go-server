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

	userID := contextKeys2.GetUserID(ctx)
	deviceID := contextKeys2.GetDeviceID(ctx)
	if userID != nil {
		userInfo.UserID = userID
	}
	if deviceID != nil {
		userInfo.DeviceID = deviceID
	}

	userInfo.RequestID = contextKeys.GetRequestID(ctx)

	return &userInfo
}
