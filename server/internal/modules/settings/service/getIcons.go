package service

import (
	"context"

	settingsModel "server/internal/modules/settings/model"
)

func (s *SettingsService) GetIcons(ctx context.Context) ([]settingsModel.Icon, error) {
	ctx, span := tracer.Start(ctx, "GetIcons")
	defer span.End()

	return s.settingsRepository.GetIcons(ctx)
}
