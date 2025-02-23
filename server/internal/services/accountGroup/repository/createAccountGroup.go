package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/services/accountGroup/repository/accountGroupDDL"
	accountGroupRepoModel "server/internal/services/accountGroup/repository/model"
)

// CreateAccountGroup создает новую группу счетов
func (r *AccountGroupRepository) CreateAccountGroup(ctx context.Context, accountGroup accountGroupRepoModel.CreateAccountGroupReq) (id uint32, serialNumber uint32, err error) {
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
		return id, serialNumber, err
	}

	// Сканируем результат
	if err = row.Scan(ctx, &serialNumber); err != nil {
		return id, serialNumber, err
	}

	// Увеличиваем серийный номер для нового элемента
	serialNumber++

	// Создаем группу счетов
	id, err = r.db.ExecWithLastInsertID(ctx, sq.
		Insert(accountGroupDDL.TableName).
		SetMap(map[string]any{
			accountGroupDDL.ColumnName:            accountGroup.Name,
			accountGroupDDL.ColumnCurrency:        accountGroup.Currency,
			accountGroupDDL.ColumnVisible:         accountGroup.Visible,
			accountGroupDDL.ColumnDatetimeCreate:  accountGroup.DatetimeCreate,
			accountGroupDDL.ColumnSerialNumber:    serialNumber,
			accountGroupDDL.ColumnCreatedByUserID: accountGroup.UserID,
		}),
	)
	if err != nil {
		return id, serialNumber, err
	}

	return id, serialNumber, nil
}
