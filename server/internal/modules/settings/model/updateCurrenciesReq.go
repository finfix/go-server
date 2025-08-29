package model

import (
	"server/internal/utils/necessary"
	"server/internal/utils/errors"
	"github.com/finfix/go-server-grpc/proto"
)

type UpdateCurrenciesReq struct {
	Necessary necessary.NecessaryUserInformation
}

// ProtoUpdateCurrenciesReq wrapper for proto request
type ProtoUpdateCurrenciesReq struct {
	*proto.UpdateCurrenciesRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoUpdateCurrenciesReq) ConvertToModel() (UpdateCurrenciesReq, error) {
	if p.UpdateCurrenciesRequest == nil {
		return UpdateCurrenciesReq{}, errors.BadRequest.New("UpdateCurrenciesRequest is required")
	}

	return UpdateCurrenciesReq{}, nil
}
