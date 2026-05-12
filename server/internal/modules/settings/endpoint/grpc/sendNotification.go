package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/settings/model"
	"server/internal/utils/necessary"

	proto "github.com/finfix/go-server-grpc/proto"
)

// SendNotification отправка уведомления пользователю
func (e *SettingsEndpoint) SendNotification(ctx context.Context, r *proto.SendNotificationRequest) (*proto.SendNotificationResponse, error) {
	res := new(proto.SendNotificationResponse)

	// Convert proto request to internal model
	req, err := model.ProtoSendNotificationReq{SendNotificationRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Parse necessary information from context
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return res, err
	}

	// Call service method
	if _, err := e.settingsService.SendNotification(ctx, req); err != nil {
		return res, err
	}

	return res, nil
}
