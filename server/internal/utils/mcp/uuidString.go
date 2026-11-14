package mcp

import (
	"github.com/google/uuid"
)

// UUIDString — указатель на UUID строкой для ответов, nil остается nil
func UUIDString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	s := id.String()
	return &s
}
