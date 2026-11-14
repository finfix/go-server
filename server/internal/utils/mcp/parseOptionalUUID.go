package mcp

import (
	"github.com/google/uuid"
)

// ParseOptionalUUID разбирает необязательный идентификатор: nil-аргумент — nil, пустая строка —
// uuid.Nil (сброс значения там, где сервис так его трактует)
func ParseOptionalUUID(field string, value *string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}
	if *value == "" {
		return &uuid.Nil, nil
	}
	id, err := ParseUUID(field, *value)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
