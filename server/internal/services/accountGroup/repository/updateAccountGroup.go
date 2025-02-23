package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/errors"

	"server/internal/services/accountGroup/model"
)

// UpdateAccountGroup обновляет группу счетов
func (r *AccountGroupRepository) UpdateAccountGroup(ctx context.Context, fields model.UpdateAccountGroupReq) error {
	ctx, span := tracer.Start(ctx, "UpdateAccountGroup")
	defer span.End()

	updates := make(map[string]any)

	// Добавляем в запрос только те поля, которые необходимо обновить
	if fields.Name != nil {
		updates["name"] = *fields.Name
	}
	if fields.Currency != nil {
		updates["currency_signatura"] = *fields.Currency
	}
	if fields.SerialNumber != nil {
		updates["serial_number"] = *fields.SerialNumber
	}
	if fields.Visible != nil {
		updates["visible"] = *fields.Visible
	}

	// Проверяем, что хоть одно поле было передано
	if len(updates) == 0 {
		return errors.BadRequest.New(ctx, "No fields to update")
	}

	// Обновляем группы счетов
	return r.db.Exec(ctx, sq.
		Update("coin.account_groups").
		SetMap(updates).
		Where(sq.Eq{"id": fields.ID}),
	)
}
