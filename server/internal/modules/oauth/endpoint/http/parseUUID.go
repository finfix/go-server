package http

import (
	"github.com/google/uuid"
)

// parseUUID разбирает client_id, пустая строка — uuid.Nil
func parseUUID(value string) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil
	}
	return id
}
