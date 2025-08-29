package model

import (
	"github.com/google/uuid"

	"pkg/datetime"
)

type Tag struct {
	ID             uuid.UUID     `json:"id" db:"id" minimum:"1"`               // Идентификатор подкатегории
	AccountGroupID uuid.UUID     `json:"accountGroupID" db:"account_group_id"` // Идентификатор группы счетов
	Name           string        `json:"name" db:"name"`                       // Название подкатегории
	DatetimeCreate datetime.Time `json:"datetimeCreate" db:"datetime_create"`  // Дата и время создания
}
