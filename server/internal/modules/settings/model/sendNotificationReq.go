package model

import (
	"github.com/google/uuid"

	userModel "server/internal/modules/user/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	"github.com/finfix/go-server-grpc/proto"
)

type SendNotificationReq struct {
	Necessary    necessary.NecessaryUserInformation
	UserID       uuid.UUID              `json:"userID"`
	Notification userModel.Notification `json:"notification"`
}

// ProtoSendNotificationReq wrapper for proto request
type ProtoSendNotificationReq struct {
	*proto.SendNotificationRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoSendNotificationReq) ConvertToModel() (SendNotificationReq, error) {
	if p.SendNotificationRequest == nil {
		return SendNotificationReq{}, errors.BadRequest.New("SendNotificationRequest is required")
	}

	// Parse UserID
	userID, err := uuid.FromBytes(p.UserID)
	if err != nil {
		return SendNotificationReq{}, errors.BadRequest.Wrap(err)
	}

	// Convert notification
	var notification userModel.Notification
	if p.Notification != nil {
		notification = userModel.Notification{
			Title:      p.Notification.Title,
			Subtitle:   p.Notification.Subtitle,
			Message:    p.Notification.Message,
			BadgeCount: uint8(p.Notification.BadgeCount),
		}
	}

	return SendNotificationReq{
		UserID:       userID,
		Notification: notification,
	}, nil
}
