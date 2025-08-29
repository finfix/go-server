package model

import (
	"server/internal/utils/errors"

	"server/internal/utils/necessary"
	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"
)

type DeleteAccountReq struct {
	Necessary necessary.NecessaryUserInformation
	ID uuid.UUID `json:"id" schema:"id" validate:"required" minimum:"1"` // Идентификатор счета
}

// ProtoDeleteAccountReq wrapper for proto request
type ProtoDeleteAccountReq struct {
	*proto.DeleteAccountRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoDeleteAccountReq) ConvertToModel() (DeleteAccountReq, error) {
	var res DeleteAccountReq

	if p.DeleteAccountRequest == nil {
		return res, errors.BadRequest.New("DeleteAccountRequest is required")
	}

	// Parse ID from bytes
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.New("Invalid ID format")
	}

	return DeleteAccountReq{
		ID: id,
	}, nil
}
