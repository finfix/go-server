package model

import (
	userModel "server/internal/services/user/model"
	"server/internal/utils/necessary"
)

type SendNotificationReq struct {
	Necessary    necessary.NecessaryUserInformation
	UserID       uint32                 `json:"userID"`
	Notification userModel.Notification `json:"notification"`
}
