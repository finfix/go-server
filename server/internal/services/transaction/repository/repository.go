package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/transaction/repository")

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB, ) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}
