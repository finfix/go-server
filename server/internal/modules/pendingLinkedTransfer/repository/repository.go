package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/modules/pendingLinkedTransfer/repository")

type PendingLinkedTransferRepository struct {
	db *sql.DB
}

func NewPendingLinkedTransferRepository(db *sql.DB) *PendingLinkedTransferRepository {
	return &PendingLinkedTransferRepository{
		db: db,
	}
}
