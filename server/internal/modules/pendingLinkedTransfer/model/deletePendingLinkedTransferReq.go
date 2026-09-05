package model

import (
	"github.com/google/uuid"

	"github.com/finfix/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type DeletePendingLinkedTransferReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"` // Идентификатор переноса
}

// ProtoDeletePendingLinkedTransferReq wrapper for proto request
type ProtoDeletePendingLinkedTransferReq struct {
	*proto.DeletePendingLinkedTransferRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoDeletePendingLinkedTransferReq) ConvertToModel() (DeletePendingLinkedTransferReq, error) {
	var res DeletePendingLinkedTransferReq

	if p.DeletePendingLinkedTransferRequest == nil {
		return res, errors.BadRequest.New("DeletePendingLinkedTransferRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	return DeletePendingLinkedTransferReq{
		ID: id,
	}, nil
}
