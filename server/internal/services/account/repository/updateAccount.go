package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/errors"
	"server/internal/services/account/repository/accountDDL"
	accountRepoModel "server/internal/services/account/repository/model"
)

// UpdateAccount обновляет счет
func (r *AccountRepository) UpdateAccount(ctx context.Context, updateReqs map[uint32]accountRepoModel.UpdateAccountReq) error {
	ctx, span := tracer.Start(ctx, "updateAccount")
	defer span.End()

	for id, fields := range updateReqs {

		// Поля для обновления
		updates := make(map[string]any)

		// Добавляем в запрос только те поля, которые необходимо обновить
		if fields.IconID != nil {
			updates[accountDDL.ColumnIconID] = *fields.IconID
		}
		if fields.AccountingInHeader != nil {
			updates[accountDDL.ColumnAccountingInHeader] = *fields.AccountingInHeader
		}
		if fields.AccountingInCharts != nil {
			updates[accountDDL.ColumnAccountingInCharts] = *fields.AccountingInCharts
		}
		if fields.Name != nil {
			updates[accountDDL.ColumnName] = *fields.Name
		}
		if fields.Visible != nil {
			updates[accountDDL.ColumnVisible] = *fields.Visible
		}
		if fields.Budget.DaysOffset != nil {
			updates[accountDDL.ColumnBudgetDaysOffset] = *fields.Budget.DaysOffset
		}
		if fields.Budget.Amount != nil {
			updates[accountDDL.ColumnBudgetAmount] = *fields.Budget.Amount
		}
		if fields.Budget.FixedSum != nil {
			updates[accountDDL.ColumnBudgetFixedSum] = *fields.Budget.FixedSum
		}
		if fields.Budget.GradualFilling != nil {
			updates[accountDDL.ColumnBudgetGradualFilling] = *fields.Budget.GradualFilling
		}
		if fields.Currency != nil {
			updates[accountDDL.ColumnCurrency] = *fields.Currency
		}
		if fields.SerialNumber != nil {
			updates[accountDDL.ColumnSerialNumber] = *fields.SerialNumber
		}
		if fields.ParentAccountID != nil {
			if *fields.ParentAccountID == 0 {
				updates[accountDDL.ColumnParentAccountID] = nil
			} else {
				updates[accountDDL.ColumnParentAccountID] = *fields.ParentAccountID
			}
		}

		// Проверяем, переданы ли поля для обновления
		if len(updates) == 0 {
			if fields.Remainder == nil {
				return errors.BadRequest.New("No fields to update")
			}
			return nil
		}

		// Обновляем счет
		err := r.db.Exec(ctx, sq.
			Update(accountDDL.Table).
			SetMap(updates).
			Where(sq.Eq{accountDDL.ColumnID: id}),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
