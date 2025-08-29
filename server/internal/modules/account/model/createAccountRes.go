package model

import (
	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"
)

type CreateAccountRes struct {
	SerialNumber                 uint32     `json:"serialNumber"`                                     // Порядковый номер счета
	BalancingTransactionID       *uuid.UUID `json:"balancingTransactionID" validate:"required"`       // Идентификатор транзакции балансировки
	BalancingAccountID           *uuid.UUID `json:"balancingAccountID" validate:"required"`           // Идентификатор балансировочного счета
	BalancingAccountSerialNumber *uint32    `json:"balancingAccountSerialNumber" validate:"required"` // Порядковый номер балансировочного счета
}

// ConvertToProto converts internal response to proto response
func (r CreateAccountRes) ConvertToProto() *proto.CreateAccountResponse {

	var balancingAccountID []byte
	if r.BalancingAccountID != nil {
		balancingAccountID = r.BalancingAccountID[:]
	}

	var balancingAccountSerialNumber *uint32
	if r.BalancingAccountSerialNumber != nil {
		balancingAccountSerialNumber = r.BalancingAccountSerialNumber
	}

	var balancingTransactionID []byte
	if r.BalancingTransactionID != nil {
		balancingTransactionID = r.BalancingTransactionID[:]
	}

	return &proto.CreateAccountResponse{
		Error:                        nil,
		SerialNumber:                 &r.SerialNumber,
		BalancingAccountID:           balancingAccountID,
		BalancingAccountSerialNumber: balancingAccountSerialNumber,
		BalancingTransactionID:       balancingTransactionID,
	}
}
