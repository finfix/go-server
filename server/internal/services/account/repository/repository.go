package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/services/account/repository")

type AccountRepository struct {
	db sql.SQL
}

func NewAccountRepository(db sql.SQL, ) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}
