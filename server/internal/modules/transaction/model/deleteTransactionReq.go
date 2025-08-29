package model

import (
	"github.com/google/uuid"

	"server/internal/utils/necessary"
	"server/internal/utils/errors"
	"github.com/finfix/go-server-grpc/proto"
)

type DeleteTransactionReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required" minimum:"1"` // Идентификатор транзакции
}

// ProtoDeleteTransactionReq wrapper for proto request
type ProtoDeleteTransactionReq struct {
	*proto.DeleteTransactionRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoDeleteTransactionReq) ConvertToModel() (DeleteTransactionReq, error) {
	var res DeleteTransactionReq

	if p.DeleteTransactionRequest == nil {
		return res, errors.BadRequest.New("DeleteTransactionRequest is required")
	}

	// Parse ID from bytes
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	return DeleteTransactionReq{
		ID: id,
	}, nil
}
