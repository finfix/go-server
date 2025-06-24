package service

import (
	"context"

	settingsModel "server/internal/services/settings/model"
)

func (s *SettingsService) SendNotification(ctx context.Context, req settingsModel.SendNotificationReq) (res settingsModel.SendNotificationRes, err error) {
	ctx, span := tracer.Start(ctx, "SendNotification")
	defer span.End()

	// Проверяем, что пользователь администратор
	err = s.checkAdmin(ctx, req.Necessary.UserID)
	if err != nil {
		return res, err
	}

	// Отправляем уведомления
	res.NotificationsSent, err = s.userService.SendNotification(
		ctx,
		req.UserID,
		req.Notification,
	)
	if err != nil {
		return res, err
	}

	return res, err
}
