package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/transaction/repository")

type TransactionRepository struct {
	db sql.SQL
}

func NewTransactionRepository(db sql.SQL, ) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}
