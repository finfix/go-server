package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"pkg/ddlHelper"
	"server/internal/modules/accountBudget/model"
	"server/internal/modules/accountBudget/repository/accountBudgetDDL"
)

// GetAccountBudgetsReq - фильтры для получения всех версий бюджета по группам счетов
type GetAccountBudgetsReq struct {
	AccountGroupIDs []uuid.UUID
	DateFrom        *time.Time
	DateTo          *time.Time
}

// GetAccountBudgets возвращает все версии бюджета по всем счетам переданных групп счетов, без объединения в "актуальную" версию
func (r *AccountBudgetRepository) GetAccountBudgets(ctx context.Context, req GetAccountBudgetsReq) (budgets []model.AccountBudget, err error) {
	ctx, span := tracer.Start(ctx, "GetAccountBudgets")
	defer span.End()

	if len(req.AccountGroupIDs) == 0 {
		return budgets, nil
	}

	q := sq.
		Select(ddlHelper.SelectAll).
		From(accountBudgetDDL.Table).
		Where(sq.Eq{
			accountBudgetDDL.ColumnAccountGroupID: req.AccountGroupIDs,
			accountBudgetDDL.ColumnIsDeleted:      false,
		}).
		OrderBy(
			accountBudgetDDL.ColumnAccountID,
			ddlHelper.Desc(accountBudgetDDL.ColumnEffectiveFrom),
		)

	if req.DateFrom != nil {
		q = q.Where(sq.GtOrEq{accountBudgetDDL.ColumnEffectiveFrom: req.DateFrom})
	}
	if req.DateTo != nil {
		q = q.Where(sq.LtOrEq{accountBudgetDDL.ColumnEffectiveFrom: req.DateTo})
	}

	return budgets, r.db.Select(ctx, &budgets, q)
}
