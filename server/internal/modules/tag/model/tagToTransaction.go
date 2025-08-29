package model

import (
	"github.com/google/uuid"
	"github.com/finfix/go-server-grpc/proto"
)

type TagToTransaction struct {
	TagID         uuid.UUID `json:"tagID" minimum:"1" db:"tag_id"`                  // Идентификатор подкатегории
	TransactionID uuid.UUID `json:"transactionID"  minimum:"1" db:"transaction_id"` // Идентификатор транзакции
}

// ConvertToProto converts TagToTransaction to proto format
func (t TagToTransaction) ConvertToProto() (*proto.TagToTransaction, error) {
	return &proto.TagToTransaction{
		TagID:         t.TagID[:],
		TransactionID: t.TransactionID[:],
	}, nil
}
