package repository

import (
	"time"

	"go.opentelemetry.io/otel"

	"pkg/cache"
	"pkg/sql"

	"server/internal/services/accountPermissions/model"
)

var tracer = otel.Tracer("/server/internal/services/accountPermissions/repository")

type AccountPermissionsRepository struct {
	db    *sql.DB
	cache *cache.Cache[struct{}, model.PermissionSet]
}

func NewAccountPermissionsRepository(db *sql.DB) *AccountPermissionsRepository {
	return &AccountPermissionsRepository{
		db:    db,
		cache: cache.NewCache[struct{}, model.PermissionSet](time.Minute),
	}
}
