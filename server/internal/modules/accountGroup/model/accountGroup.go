package model

import (
	"github.com/google/uuid"

	"pkg/datetime"
)

type AccountGroup struct {
	ID             uuid.UUID     `json:"id" db:"id"`                          // Идентификатор группы счетов
	Name           string        `json:"name" db:"name"`                      // Название группы счетов
	Currency       string        `json:"currency" db:"currency_signatura"`    // Валюта группы счетов
	SerialNumber   uuid.UUID     `json:"serialNumber" db:"serial_number"`     // Порядковый номер группы счетов
	Visible        bool          `json:"visible" db:"visible"`                // Видимость группы счетов
	DatetimeCreate datetime.Time `json:"datetimeCreate" db:"datetime_create"` // Дата и время создания группы счетов
}
