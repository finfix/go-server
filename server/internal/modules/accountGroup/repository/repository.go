package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/modules/accountGroup/repository")

type AccountGroupRepository struct {
	db *sql.DB
}

func NewAccountGroupRepository(db *sql.DB) *AccountGroupRepository {
	return &AccountGroupRepository{
		db: db,
	}
}
