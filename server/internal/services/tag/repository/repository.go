package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/services/tag/repository")

type TagRepository struct {
	db sql.SQL
}

func NewTagRepository(db sql.SQL, ) *TagRepository {
	return &TagRepository{
		db: db,
	}
}
