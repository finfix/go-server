package service

import (
	"context"

	"pkg/errors"

	settingsModel "server/internal/services/settings/model"
	"server/internal/services/settings/model/applicationType"
)

func (s *SettingsService) GetVersion(ctx context.Context, appType applicationType.Type) (version settingsModel.Version, err error) {
	ctx, span := tracer.Start(ctx, "GetVersion")
	defer span.End()

	switch appType {
	case applicationType.Server:
		return settingsModel.Version{
			Version: s.version.Version,
			Build:   s.version.Build,
		}, nil
	case applicationType.IOs:
		return s.settingsRepository.GetVersion(ctx, appType)
	case applicationType.Android, applicationType.Web:
		return version, errors.NotFound.New(ctx, "Такое приложение еще не реализовано")
	default:
		return version, errors.BadRequest.New(ctx, "Неверный тип приложения")
	}
}
