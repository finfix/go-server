package model

import "github.com/google/uuid"

type CreateAccountRes struct {
	ID                           uuid.UUID  `json:"id"`                                               // Идентификатор созданного счета
	SerialNumber                 uint32     `json:"serialNumber"`                                     // Порядковый номер счета
	BalancingTransactionID       *uuid.UUID `json:"balancingTransactionID" validate:"required"`       // Идентификатор транзакции балансировки
	BalancingAccountID           *uuid.UUID `json:"balancingAccountID" validate:"required"`           // Идентификатор балансировочного счета
	BalancingAccountSerialNumber *uint32    `json:"balancingAccountSerialNumber" validate:"required"` // Порядковый номер балансировочного счета
}
