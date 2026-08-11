package model

import "github.com/google/uuid"

type Device struct {
	ID                     uuid.UUID  `db:"id" json:"id"` // Идентификатор девайса
	DeviceInformation                 // Информация о девайсе пользователя
	ApplicationInformation            // Информация о приложении пользователя
	NotificationToken      *string    `db:"notification_token" json:"-"`       // Токен для системы уведомлений
	RefreshToken           string     `db:"refresh_token" json:"-"`            // Токен доступа для обновления пары токенов
	UserID                 uuid.UUID  `db:"user_id" json:"userID"`             // Идентификатор пользователя девайса
	DeviceID               string     `db:"device_id" json:"deviceID"`         // Идентификатор девайса
	LastSyncedAuditLogID   uint32     `db:"last_synced_audit_log_id" json:"-"` // Идентификатор последней подтвержденной (Confirm) записи аудит-лога
	PendingCheckpoint      *uint32    `db:"pending_checkpoint" json:"-"`       // Чекпоинт последнего Sync, ожидающий подтверждения
	PendingSyncToken       *uuid.UUID `db:"pending_sync_token" json:"-"`       // Опорный токен последнего Sync
}
