package auth

import (
	"context"

	"server/internal/utils/errors"
)

type Claims struct {
	UserID   uint32
	DeviceID string
}

func (c *Claims) Validate(ctx context.Context) error {
	if c.UserID == 0 {
		return errors.Unauthorized.New("UserID is empty").WithContextParams(ctx)
	}
	if c.DeviceID == "" {
		return errors.Unauthorized.New("DeviceID is empty").WithContextParams(ctx)
	}
	return nil
}
