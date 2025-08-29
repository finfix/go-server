package model

import (
	"github.com/finfix/go-server-grpc/proto"
)

type CreateTransactionRes struct {
}

// ConvertToProto converts CreateTransactionRes to proto format
func (r CreateTransactionRes) ConvertToProto() *proto.CreateTransactionResponse {
	return &proto.CreateTransactionResponse{}
}

type GetTransactionsRes struct {
	Transactions []Transaction `json:"transactions"`
}
