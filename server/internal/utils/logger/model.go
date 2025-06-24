package logger

type UserInfo struct {
	UserID    *uint32 `json:"userID,omitempty"`
	RequestID *string `json:"requestID,omitempty"`
	DeviceID  *string `json:"deviceID,omitempty"`
}
