package model

import (
	"github.com/google/uuid"

	"server/internal/utils/necessary"
	"server/internal/utils/errors"
	"github.com/finfix/go-server-grpc/proto"
)

type DeleteAccountGroupReq struct {
	Necessary necessary.NecessaryUserInformation
	ID uuid.UUID `json:"id" schema:"id" validate:"required" minimum:"1"` // Идентификатор счета
}

// ProtoDeleteAccountGroupReq wrapper for proto request
type ProtoDeleteAccountGroupReq struct {
	*proto.DeleteAccountGroupRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoDeleteAccountGroupReq) ConvertToModel() (DeleteAccountGroupReq, error) {
	var res DeleteAccountGroupReq

	if p.DeleteAccountGroupRequest == nil {
		return res, errors.BadRequest.New("DeleteAccountGroupRequest is required")
	}

	// Parse ID from bytes
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	return DeleteAccountGroupReq{
		ID: id,
	}, nil
}
