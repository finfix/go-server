package model

import "github.com/google/uuid"

type UpdateAccountRes struct {
	BalancingAccountID           *uuid.UUID `json:"balancingAccountID" validate:"required"`           // Идентификатор балансировочного счета
	BalancingTransactionID       *uuid.UUID `json:"balancingTransactionID" validate:"required"`       // Идентификатор транзакции
	BalancingAccountSerialNumber *uint32    `json:"balancingAccountSerialNumber" validate:"required"` // Порядковый номер балансировочного счета
}
