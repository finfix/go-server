package model

import (
	"github.com/google/uuid"

	userModel "server/internal/modules/user/model"
	"server/internal/utils/necessary"
)

type SendNotificationReq struct {
	Necessary    necessary.NecessaryUserInformation
	UserID       uuid.UUID              `json:"userID"`
	Notification userModel.Notification `json:"notification"`
}
