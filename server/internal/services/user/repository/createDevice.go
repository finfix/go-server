package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	userModel "server/internal/services/user/model"
	"server/internal/services/user/repository/deviceDDL"
)

// CreateDevice Создает новый девайс для пользователя
func (r *UserRepository) CreateDevice(ctx context.Context, req userModel.Device) (id uint32, err error) {
	ctx, span := tracer.Start(ctx, "CreateDevice")
	defer span.End()

	return r.db.ExecWithLastInsertID(ctx, sq.
		Insert(deviceDDL.Table).
		SetMap(map[string]any{
			deviceDDL.ColumnRefreshToken:        req.RefreshToken,
			deviceDDL.ColumnDeviceID:            req.DeviceID,
			deviceDDL.ColumnUserID:              req.UserID,
			deviceDDL.ColumnDeviceOSName:        req.DeviceInformation.NameOS,
			deviceDDL.ColumnDeviceOSVersion:     req.DeviceInformation.VersionOS,
			deviceDDL.ColumnDeviceName:          req.DeviceInformation.DeviceName,
			deviceDDL.ColumnDeviceModelName:     req.DeviceInformation.ModelName,
			deviceDDL.ColumnDeviceIPAddress:     req.DeviceInformation.IPAddress,
			deviceDDL.ColumnDeviceUserAgent:     req.DeviceInformation.UserAgent,
			deviceDDL.ColumnApplicationBundleID: req.ApplicationInformation.BundleID,
			deviceDDL.ColumnApplicationVersion:  req.ApplicationInformation.Version,
			deviceDDL.ColumnApplicationBuild:    req.ApplicationInformation.Build,
		}),
	)
}
