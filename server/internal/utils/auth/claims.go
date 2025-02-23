package auth

import (
	"context"

	"pkg/errors"
)

type Claims struct {
	UserID   uint32
	DeviceID string
}

func (c *Claims) Validate(ctx context.Context) error {
	if c.UserID == 0 {
		return errors.Unauthorized.New(ctx, "UserID is empty")
	}
	if c.DeviceID == "" {
		return errors.Unauthorized.New(ctx, "DeviceID is empty")
	}
	return nil
}
