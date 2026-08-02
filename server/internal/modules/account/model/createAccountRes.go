package model

import (
	"github.com/finfix/go-server-grpc/proto"
)

type CreateAccountRes struct {
	SerialNumber uint32 `json:"serialNumber"` // Порядковый номер счета
}

// ConvertToProto converts internal response to proto response
func (r CreateAccountRes) ConvertToProto() *proto.CreateAccountResponse {
	return &proto.CreateAccountResponse{
		Error:        nil,
		SerialNumber: &r.SerialNumber,
	}
}
