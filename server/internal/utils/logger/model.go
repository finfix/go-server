package logger

import "github.com/google/uuid"

type UserInfo struct {
	UserID    *uuid.UUID `json:"userID,omitempty"`
	RequestID *string `json:"requestID,omitempty"`
	DeviceID  *string `json:"deviceID,omitempty"`
}
