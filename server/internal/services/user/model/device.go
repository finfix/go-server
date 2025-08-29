package model

import "github.com/google/uuid"

type Device struct {
	DeviceInformation              // Информация о девайсе пользователя
	ApplicationInformation         // Информация о приложении пользователя
	NotificationToken      *string `db:"notification_token" json:"-"` // Токен для системы уведомлений
	RefreshToken           string  `db:"refresh_token" json:"-"` // Токен доступа для обновления пары токенов
	UserID                 uuid.UUID  `db:"user_id" json:"userID"` // Идентификатор пользователя девайса
	DeviceID               string  `db:"device_id" json:"deviceID"` // Идентификатор девайса
}
