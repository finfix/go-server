package repository

import (
	"go.opentelemetry.io/otel"

	"github.com/google/uuid"

	"pkg/cache"
	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/modules/userToAccountGroup/repository")

// UserToAccountGroupRepository - репозиторий связей пользователей с группами счетов
type UserToAccountGroupRepository struct {
	db                           *sql.DB
	accessedAccountGroupIDsCache *cache.ItemCache[uuid.UUID, []uuid.UUID] // Кэш юзер - массив доступных ему групп счетов
}

// NewUserToAccountGroupRepository создает новый репозиторий связей пользователей с группами счетов
func NewUserToAccountGroupRepository(db *sql.DB) *UserToAccountGroupRepository {
	return &UserToAccountGroupRepository{
		db:                           db,
		accessedAccountGroupIDsCache: cache.NewItemCache[uuid.UUID, []uuid.UUID](),
	}
}
