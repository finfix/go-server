package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/services/user/repository/deviceDDL"
)

// DeleteDevice Удаляет девайс пользователя
func (r *UserRepository) DeleteDevice(ctx context.Context, userID uint32, deviceID string) error {
	ctx, span := tracer.Start(ctx, "DeleteDevice")
	defer span.End()

	return r.db.Exec(ctx, sq.
		Delete(deviceDDL.Table).
		Where(sq.Eq{
			deviceDDL.ColumnUserID:   userID,
			deviceDDL.ColumnDeviceID: deviceID,
		}),
	)
}
