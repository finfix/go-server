package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/services/account/repository/accountDDL"
	accountRepoModel "server/internal/services/account/repository/model"
)

// CreateAccount создает новый счет
func (r *AccountRepository) CreateAccount(ctx context.Context, account accountRepoModel.CreateAccountReq) (id uint32, serialNumber uint32, err error) {
	ctx, span := tracer.Start(ctx, "createAccount")
	defer span.End()

	// Получаем максимальный серийный номер в группе счетов
	row, err := r.db.QueryRow(ctx, sq.
		Select(ddlHelper.Coalesce(
			ddlHelper.Max(accountDDL.ColumnSerialNumber),
			"1",
		)).
		From(accountDDL.Table).
		Where(sq.Eq{accountDDL.ColumnAccountGroupID: account.AccountGroupID}),
	)
	if err != nil {
		return id, serialNumber, err
	}

	// Сканируем результат
	if err = row.Scan(ctx, &serialNumber); err != nil {
		return id, serialNumber, err
	}

	// Увеличиваем серийный номер для нового элемента
	serialNumber++

	// Создаем счет
	id, err = r.db.ExecWithLastInsertID(ctx, sq.
		Insert(accountDDL.Table).
		SetMap(map[string]any{
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
			accountDDL.ColumnSerialNumber:         serialNumber,
		}))
	if err != nil {
		return id, serialNumber, err
	}

	return id, serialNumber, nil
}
