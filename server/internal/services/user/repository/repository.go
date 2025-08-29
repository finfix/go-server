package repository

import (
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	"pkg/cache"
	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/services/user/repository")

type UserRepository struct {
	db                           *sql.DB
	accessedAccountGroupIDsCache *cache.ItemCache[uuid.UUID, []uuid.UUID] // Кэш юзер - массив доступных ему групп счетов
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db:                           db,
		accessedAccountGroupIDsCache: cache.NewItemCache[uuid.UUID, []uuid.UUID](time.Minute),
	}
}
