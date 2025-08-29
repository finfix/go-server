package service

import (
	"context"
	"server/internal/enum/applicationType"

	"server/internal/utils/errors"

	settingsModel "server/internal/modules/settings/model"
)

func (s *SettingsService) GetVersion(ctx context.Context, req settingsModel.GetVersionReq) (res settingsModel.GetVersionRes, err error) {
	ctx, span := tracer.Start(ctx, "GetVersion")
	defer span.End()

	switch req.ApplicationType {
	case applicationType.Server:

		return settingsModel.GetVersionRes{
			Version: settingsModel.Version{
				Version: s.version.Version,
				Build:   s.version.Build,
			},
		}, nil

	case applicationType.IOS:
		version, err := s.settingsRepository.GetVersion(ctx, req)
		if err != nil {
			return res, err
		}

		return settingsModel.GetVersionRes{Version: version}, nil
	case applicationType.Android, applicationType.Web:
		return res, errors.NotFound.New("Такое приложение еще не реализовано").WithContextParams(ctx)
	default:
		return res, errors.BadRequest.New("Неверный тип приложения").WithContextParams(ctx)
	}
}
