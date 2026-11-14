package mcp

import (
	"fmt"

	"github.com/google/uuid"
)

// ParseUUID разбирает обязательный идентификатор из аргумента инструмента
func ParseUUID(field, value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: ожидается UUID, получено %q", field, value)
	}
	return id, nil
}
