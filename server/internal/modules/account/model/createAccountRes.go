package model

import (
	"github.com/finfix/go-server-grpc/proto"
)

type CreateAccountRes struct{}

// ConvertToProto converts internal response to proto response
func (r CreateAccountRes) ConvertToProto() *proto.CreateAccountResponse {
	return &proto.CreateAccountResponse{
		Error: nil,
	}
}
