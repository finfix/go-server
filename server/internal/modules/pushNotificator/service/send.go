package service

import (
	"context"
	"time"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"

	"pkg/log"
	"server/internal/utils/errors"

	"server/internal/modules/pushNotificator/model"
)

// SendNotification отправляет одно сообщение на все переданные устройства
func (s *PushNotificatorService) SendNotification(ctx context.Context, req model.SendNotificationReq) (id string, err error) {
	ctx, span := tracer.Start(ctx, "SendNotification")
	defer span.End()

	const defaultPriority = 5

	if !s.isOn {
		log.WithContextParams(ctx).Warning("Вызвана функция SendNotification. Пуши выключены")
		return id, nil
	}

	payload := payload.NewPayload()
	if req.Notification.Title != nil {
		payload = payload.AlertTitle(*req.Notification.Title)
	}
	if req.Notification.Subtitle != nil {
		payload = payload.AlertSubtitle(*req.Notification.Subtitle)
	}
	if req.Notification.Message != nil {
		payload = payload.AlertBody(*req.Notification.Message)
	}
	if req.Notification.Badge != nil {
		payload = payload.Badge(int(*req.Notification.Badge))
	}
	payload = payload.ContentAvailable()

	notification := &apns2.Notification{
		ApnsID:      "",
		CollapseID:  "",
		DeviceToken: req.NotificationToken,
		Topic:       req.BundleID,
		Expiration:  time.Time{},
		Priority:    defaultPriority,
		Payload:     payload,
		PushType:    apns2.PushTypeAlert,
	}

	res, err := s.apns.PushWithContext(ctx, notification)
	if err != nil {
		return id, errors.InternalServer.Wrap(err).WithContextParams(ctx)
	}
	id = res.ApnsID

	if !res.Sent() {
		return id, errors.InternalServer.New(res.Reason).WithContextParams(ctx)
	}

	return id, nil
}
