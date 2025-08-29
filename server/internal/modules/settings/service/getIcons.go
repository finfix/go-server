package service

import (
	"context"

	settingsModel "server/internal/modules/settings/model"
)

func (s *SettingsService) GetIcons(ctx context.Context) (res settingsModel.GetIconsRes, err error) {
	ctx, span := tracer.Start(ctx, "GetIcons")
	defer span.End()

	icons, err := s.settingsRepository.GetIcons(ctx)
	if err != nil {
		return res, err
	}

	return settingsModel.GetIconsRes{
		Icons: icons,
	}, nil
}
