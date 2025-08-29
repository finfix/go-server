package service

import (
	"context"

	"github.com/google/uuid"
	"pkg/log"

	model2 "server/internal/modules/pushNotificator/model"
	"server/internal/modules/user/model"
	"server/internal/modules/user/model/OS"
	userRepoModel "server/internal/modules/user/repository/model"
)

// SendNotification отправляет пуш на все устройства пользователя
func (s *UserService) SendNotification(ctx context.Context, userID uuid.UUID, push model.Notification) (count uint8, err error) {
	ctx, span := tracer.Start(ctx, "SendNotification")
	defer span.End()

	// Получаем все девайсы пользователя
	devices, err := s.userRepository.GetDevices(ctx, userRepoModel.GetDevicesReq{ //nolint:exhaustruct
		UserIDs: []uuid.UUID{userID},
	})
	if err != nil {
		return count, err
	}

	// Проходимся по каждому девайсу
	for _, device := range devices {

		if device.NotificationToken == nil {
			continue
		}

		// Смотрим на операционную систему и отправляем уведомление
		switch device.DeviceInformation.NameOS {
		case OS.IOS, OS.IPadOS, OS.OSX, OS.WatchOS:
			_, err = s.pushNotificator.SendNotification(ctx, model2.SendNotificationReq{
				Notification: model2.NotificationSettings{
					Title:    &push.Title,
					Message:  &push.Message,
					Subtitle: &push.Subtitle,
					Badge:    &push.BadgeCount,
				},
				NotificationToken: *device.NotificationToken,
				BundleID:          device.ApplicationInformation.BundleID,
			})

		case OS.Android:
			break
		}
		if err != nil {
			log.WithContextParams(ctx).Error(err)
		}
		count++
	}

	return count, nil
}
