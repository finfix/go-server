package model

import "github.com/google/uuid"

type CreateAccountGroupRes struct {
	ID           uuid.UUID `json:"id"`           // Идентификатор созданной группы счетов
	SerialNumber uint32    `json:"serialNumber"` // Порядковый номер группы счетов
}
