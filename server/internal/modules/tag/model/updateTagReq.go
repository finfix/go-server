package model

import (
	"github.com/google/uuid"

	"server/internal/utils/necessary"
	"server/internal/utils/errors"
	"github.com/finfix/go-server-grpc/proto"
)

type UpdateTagReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID  `json:"id" validate:"required" minimum:"1"` // Идентификатор подкатегории
	Name      *string `json:"name" minimum:"1"`                   // Название подкатегории
}

// ProtoUpdateTagReq wrapper for proto request
type ProtoUpdateTagReq struct {
	*proto.UpdateTagRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoUpdateTagReq) ConvertToModel() (UpdateTagReq, error) {
	var res UpdateTagReq

	if p.UpdateTagRequest == nil {
		return res, errors.BadRequest.New("UpdateTagRequest is required")
	}

	// Parse ID from bytes
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	return UpdateTagReq{
		ID:   id,
		Name: p.Name,
	}, nil
}
