package model

import "github.com/google/uuid"

// UserToAccountGroup - связь пользователя с группой счетов
type UserToAccountGroup struct {
	UserID         uuid.UUID `db:"user_id"`          // Идентификатор пользователя
	AccountGroupID uuid.UUID `db:"account_group_id"` // Идентификатор группы счетов
}
