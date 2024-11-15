package repository

import (
	"time"

	"go.opentelemetry.io/otel"

	"pkg/cache"
	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/services/user/repository")

type UserRepository struct {
	db                           sql.SQL
	accessedAccountGroupIDsCache *cache.Cache[uint32, []uint32] // Кэш юзер - массив доступных ему групп счетов
}

func NewUserRepository(db sql.SQL) *UserRepository {
	return &UserRepository{
		db:                           db,
		accessedAccountGroupIDsCache: cache.NewCache[uint32, []uint32](time.Minute),
	}
}
