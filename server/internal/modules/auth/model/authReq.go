package model

import "github.com/google/uuid"

type AuthRes struct {
	Tokens `json:"token"`     // Токены доступа
	ID     uuid.UUID `json:"id"` // Идентификатор пользователя
}
