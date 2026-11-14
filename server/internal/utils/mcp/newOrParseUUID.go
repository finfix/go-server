package mcp

import (
	"github.com/google/uuid"
)

// NewOrParseUUID возвращает идентификатор из аргумента или новый, если агент его не передал
func NewOrParseUUID(field string, value *string) (uuid.UUID, error) {
	if value == nil || *value == "" {
		return uuid.New(), nil
	}
	return ParseUUID(field, *value)
}
