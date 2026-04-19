package model

import (
	"github.com/finfix/go-server-grpc/proto"
)

type CreateTagRes struct {
	// ID теперь известен заранее и не возвращается в ответе
}

// ConvertToProto converts CreateTagRes to proto format
func (r CreateTagRes) ConvertToProto() *proto.CreateTagResponse {
	return &proto.CreateTagResponse{}
}
