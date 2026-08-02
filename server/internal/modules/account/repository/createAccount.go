package repository

import (
	"context"

	"server/internal/modules/account/repository/accountDDL"
	accountRepoModel "server/internal/modules/account/repository/model"

	sq "github.com/Masterminds/squirrel"
)

// CreateAccount создает новый счет
func (r *AccountRepository) CreateAccount(ctx context.Context, account accountRepoModel.CreateAccountReq) error {
	ctx, span := tracer.Start(ctx, "createAccount")
	defer span.End()

	// Создаем счет
	return r.db.Exec(ctx, sq.
		Insert(accountDDL.Table).
		SetMap(map[string]any{
			accountDDL.ColumnID:                   account.ID,
			accountDDL.ColumnBudgetAmount:         account.Budget.Amount,
			accountDDL.ColumnName:                 account.Name,
			accountDDL.ColumnIconID:               account.IconID,
			accountDDL.ColumnType:                 account.Type,
			accountDDL.ColumnCurrency:             account.Currency,
			accountDDL.ColumnVisible:              account.Visible,
			accountDDL.ColumnAccountGroupID:       account.AccountGroupID,
			accountDDL.ColumnAccountingInHeader:   account.AccountingInHeader,
			accountDDL.ColumnAccountingInCharts:   account.AccountingInCharts,
			accountDDL.ColumnBudgetGradualFilling: account.Budget.GradualFilling,
			accountDDL.ColumnIsParent:             account.IsParent,
			accountDDL.ColumnBudgetFixedSum:       account.Budget.FixedSum,
			accountDDL.ColumnBudgetDaysOffset:     account.Budget.DaysOffset,
			accountDDL.ColumnParentAccountID:      account.ParentAccountID,
			accountDDL.ColumnCreatedByUserID:      account.UserID,
			accountDDL.ColumnDatetimeCreate:       account.DatetimeCreate,
			accountDDL.ColumnRank:                 account.Rank,
		}))
}
