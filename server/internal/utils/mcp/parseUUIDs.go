package mcp

import (
	"github.com/google/uuid"
)

// ParseUUIDs разбирает список идентификаторов
func ParseUUIDs(field string, values []string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		id, err := ParseUUID(field, value)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
