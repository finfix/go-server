package model

import "github.com/google/uuid"

type CreateRes struct {
	ID uuid.UUID `json:"id"`
}
