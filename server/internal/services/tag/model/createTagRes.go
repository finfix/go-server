package model

import "github.com/google/uuid"

type CreateTagRes struct {
	ID uuid.UUID `json:"id" validate:"required" minimum:"1"` // Идентификатор транзакции
}
