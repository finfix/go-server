package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	settingsModel "server/internal/services/settings/model"
	"server/internal/services/settings/repository/iconDDL"
)

func (r *SettingsRepository) GetIcons(ctx context.Context) (icons []settingsModel.Icon, err error) {
	ctx, span := tracer.Start(ctx, "GetIcons")
	defer span.End()

	return icons, r.db.Select(ctx, &icons, sq.
		Select(ddlHelper.SelectAll).
		From(iconDDL.Table),
	)
}
