package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/accountBudget/model"
	"server/internal/modules/accountBudget/repository/accountBudgetDDL"
)

// CreateAccountBudget создает новую версию бюджета счета
func (r *AccountBudgetRepository) CreateAccountBudget(ctx context.Context, req model.AccountBudget) error {
	ctx, span := tracer.Start(ctx, "CreateAccountBudget")
	defer span.End()

	// Создаем версию бюджета
	return r.db.Exec(ctx, sq.
		Insert(accountBudgetDDL.Table).
		SetMap(map[string]any{
			accountBudgetDDL.ColumnID:              req.ID,
			accountBudgetDDL.ColumnAccountID:       req.AccountID,
			accountBudgetDDL.ColumnAccountGroupID:  req.AccountGroupID,
			accountBudgetDDL.ColumnAmount:          req.Amount,
			accountBudgetDDL.ColumnFixedSum:        req.FixedSum,
			accountBudgetDDL.ColumnDaysOffset:      req.DaysOffset,
			accountBudgetDDL.ColumnGradualFilling:  req.GradualFilling,
			accountBudgetDDL.ColumnEffectiveFrom:   req.EffectiveFrom,
			accountBudgetDDL.ColumnCreatedByUserID: req.CreatedByUserID,
		}),
	)
}
