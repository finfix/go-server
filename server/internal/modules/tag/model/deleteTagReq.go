package model

import (
	"github.com/google/uuid"

	"server/internal/utils/necessary"
	"server/internal/utils/errors"
	"github.com/finfix/go-server-grpc/proto"
)

type DeleteTagReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required" minimum:"1"` // Идентификатор подкатегории
}

// ProtoDeleteTagReq wrapper for proto request
type ProtoDeleteTagReq struct {
	*proto.DeleteTagRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoDeleteTagReq) ConvertToModel() (DeleteTagReq, error) {
	var res DeleteTagReq

	if p.DeleteTagRequest == nil {
		return res, errors.BadRequest.New("DeleteTagRequest is required")
	}

	// Parse ID from bytes
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	return DeleteTagReq{
		ID: id,
	}, nil
}
