package repository

import (
	"go.opentelemetry.io/otel"

	"pkg/sql"
)

var tracer = otel.Tracer("/server/internal/modules/accountBudget/repository")

// AccountBudgetRepository - репозиторий версий бюджета счетов
type AccountBudgetRepository struct {
	db *sql.DB
}

// NewAccountBudgetRepository создает новый репозиторий версий бюджета счетов
func NewAccountBudgetRepository(db *sql.DB) *AccountBudgetRepository {
	return &AccountBudgetRepository{
		db: db,
	}
}
