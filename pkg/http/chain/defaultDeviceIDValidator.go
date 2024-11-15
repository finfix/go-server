package chain

import (
	"context"
	"net/http"

	"pkg/contextKeys"
	"pkg/errors"
)

func DefaultDeviceIDValidator(ctx context.Context, r *http.Request) (context.Context, error) {
	_, span := tracer.Start(ctx, "DefaultDeviceIDValidator")
	defer span.End()

	deviceID := r.Header.Get("DeviceID")
	if deviceID == "" {
		return ctx, errors.BadRequest.New("DeviceID is empty")
	}
	ctx = contextKeys.SetUserID(ctx, 0)
	return contextKeys.SetDeviceID(ctx, deviceID), nil
}
