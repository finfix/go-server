package repository

import (
	"context"

	"pkg/ddlHelper"
	"server/internal/modules/accountGroup/repository/accountGroupDDL"
	accountGroupRepoModel "server/internal/modules/accountGroup/repository/model"

	sq "github.com/Masterminds/squirrel"
)

// CreateAccountGroup создает новую группу счетов
func (r *AccountGroupRepository) CreateAccountGroup(ctx context.Context, accountGroup accountGroupRepoModel.CreateAccountGroupReq) (serialNumber uint32, err error) {
	ctx, span := tracer.Start(ctx, "CreateAccountGroup")
	defer span.End()

	// Получаем текущий максимальный серийный номер группы счетов
	row, err := r.db.QueryRow(ctx, sq.
		Select(ddlHelper.Coalesce(
			ddlHelper.Max(accountGroupDDL.ColumnSerialNumber),
			"1",
		)).
		From(accountGroupDDL.TableName),
	)
	if err != nil {
		return serialNumber, err
	}

	// Сканируем результат
	if err = row.Scan(&serialNumber); err != nil {
		return serialNumber, err
	}

	// Увеличиваем серийный номер для нового элемента
	serialNumber++

	// Создаем группу счетов
	return serialNumber, r.db.Exec(ctx, sq.
		Insert(accountGroupDDL.TableName).
		SetMap(map[string]any{
			accountGroupDDL.ColumnID:              accountGroup.ID,
			accountGroupDDL.ColumnName:            accountGroup.Name,
			accountGroupDDL.ColumnCurrency:        accountGroup.Currency,
			accountGroupDDL.ColumnVisible:         accountGroup.Visible,
			accountGroupDDL.ColumnDatetimeCreate:  accountGroup.DatetimeCreate,
			accountGroupDDL.ColumnSerialNumber:    serialNumber,
			accountGroupDDL.ColumnCreatedByUserID: accountGroup.UserID,
		}),
	)
}
