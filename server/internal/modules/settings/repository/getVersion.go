package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	settingsModel "server/internal/modules/settings/model"
	"server/internal/modules/settings/model/applicationType"
	"server/internal/modules/settings/repository/versionDDL"
)

func (r *SettingsRepository) GetVersion(ctx context.Context, appType applicationType.Type) (version settingsModel.Version, err error) {
	ctx, span := tracer.Start(ctx, "GetVersion")
	defer span.End()

	return version, r.db.Get(ctx, &version, sq.
		Select(ddlHelper.SelectAll).
		From(versionDDL.Table).
		Where(sq.Eq{versionDDL.ColumnName: appType}).
		Limit(1),
	)
}
