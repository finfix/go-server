package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/services/accountGroup/repository")

type AccountGroupRepository struct {
	db sql.SQL
}

func NewAccountGroupRepository(db sql.SQL, ) *AccountGroupRepository {
	return &AccountGroupRepository{
		db: db,
	}
}
