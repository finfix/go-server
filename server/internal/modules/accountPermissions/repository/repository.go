package repository

import (
	"time"

	"go.opentelemetry.io/otel"

	"pkg/cache"
	"pkg/sql"

	"server/internal/modules/accountPermissions/model"
)

var tracer = otel.Tracer("/server/internal/modules/accountPermissions/repository")

type AccountPermissionsRepository struct {
	db    *sql.DB
	cache *cache.ItemCache[struct{}, model.PermissionSet]
}

func NewAccountPermissionsRepository(db *sql.DB) *AccountPermissionsRepository {
	return &AccountPermissionsRepository{
		db:    db,
		cache: cache.NewItemCache[struct{}, model.PermissionSet](time.Minute),
	}
}
