package model

import "github.com/google/uuid"

type CreateTransactionRes struct {
	ID uuid.UUID `json:"id" validate:"required" minimum:"1"` // Идентификатор транзакции
}

type GetTransactionsRes struct {
	Transactions []Transaction `json:"transactions"`
}
