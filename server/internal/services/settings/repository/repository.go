package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/services/settings/repository")

type SettingsRepository struct {
	db sql.SQL
}

func NewSettingsRepository(db sql.SQL) *SettingsRepository {
	return &SettingsRepository{
		db: db,
	}
}
