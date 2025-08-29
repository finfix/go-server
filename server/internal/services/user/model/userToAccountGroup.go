package model

import "github.com/google/uuid"

type UserToAccountGroup struct {
	UserID         uuid.UUID `db:"user_id"`
	AccountGroupID uuid.UUID `db:"account_group_id"`
}
