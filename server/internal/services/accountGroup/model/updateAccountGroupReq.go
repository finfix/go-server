package model

import (
	"github.com/google/uuid"

	"server/internal/utils/necessary"
)

type UpdateAccountGroupReq struct {
	Necessary    necessary.NecessaryUserInformation
	ID           uuid.UUID  `json:"id" db:"id"`                       // Идентификатор группы счетов
	Name         *string    `json:"name" db:"name"`                   // Название группы счетов
	Currency     *string    `json:"currency" db:"currency_signatura"` // Валюта группы счетов
	Visible      *bool      `json:"visible" db:"visible"`             // Видимость группы счетов
	SerialNumber *uuid.UUID `json:"serialNumber" db:"serial_number"`  // Порядковый номер группы счетов
}
