package model

import "github.com/google/uuid"

type TagToTransaction struct {
	TagID         uuid.UUID `json:"tagID" minimum:"1" db:"tag_id"`                  // Идентификатор подкатегории
	TransactionID uuid.UUID `json:"transactionID"  minimum:"1" db:"transaction_id"` // Идентификатор транзакции
}
