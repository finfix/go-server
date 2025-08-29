package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/services/account/repository/accountDDL"
	"server/internal/utils/errors"
)

// CheckAccess проверяет, имеет ли набор групп счетов пользователя доступ к указанным идентификаторам счетов
func (r *AccountRepository) CheckAccess(ctx context.Context, accountGroupIDs, accountIDs []uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "CheckAccess")
	defer span.End()

	// Получаем все доступные счета по группам счетов и перечисленным счетам
	rows, err := r.db.Query(ctx, sq.
		Select(accountDDL.ColumnID).
		From(accountDDL.Table).
		Where(sq.Eq{
			accountDDL.ColumnAccountGroupID: accountGroupIDs,
			accountDDL.ColumnID:             accountIDs,
		}),
	)
	if err != nil {
		return err
	}

	// Формируем мапу доступных счетов
	accessedAccountIDs := make(map[uuid.UUID]struct{})

	// Проходимся по каждому доступному счету
	for rows.Next() {

		// Считываем ID счета
		var accountID uuid.UUID
		if err = rows.Scan(&accountID); err != nil {
			return err
		}

		// Добавляем ID счета в мапу
		accessedAccountIDs[accountID] = struct{}{}
	}

	if len(accessedAccountIDs) == 0 {
		return errors.Forbidden.New("You don't have access to any of the requested accounts").WithContextParams(ctx).
			WithParams(
				"AccountGroupIDs", accountGroupIDs,
				"AccountIDs", accountIDs,
			)
	}

	// Проходимся по каждому запрашиваемому счету
	for _, accountID := range accountIDs {

		// Если счета нет в мапе доступных счетов, возвращаем ошибку
		if _, ok := accessedAccountIDs[accountID]; !ok {
			return errors.Forbidden.New("You don't have access to account").
				WithContextParams(ctx).
				WithParams(
					"AccountGroupIDs", accountGroupIDs,
					"AccountID", accountID,
				).
				SkipPreviousCaller()
		}
	}

	return nil
}
