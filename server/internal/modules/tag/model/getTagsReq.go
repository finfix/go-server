package model

import (
	"github.com/google/uuid"

	"github.com/finfix/go-server-grpc/proto"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type GetTagsReq struct {
	Necessary       necessary.NecessaryUserInformation
	AccountGroupIDs []uuid.UUID `json:"accountGroupIDs" schema:"accountGroupIDs" minimum:"1"` // Идентификаторы групп счетов
	IDs             []uuid.UUID `json:"ids" schema:"ids" minimum:"1"`                         // Идентификаторы подкатегорий
}

// ProtoGetTagsReq wrapper for proto request
type ProtoGetTagsReq struct {
	*proto.GetTagsRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoGetTagsReq) ConvertToModel() (GetTagsReq, error) {
	var res GetTagsReq

	if p.GetTagsRequest == nil {
		return res, errors.BadRequest.New("GetTagsRequest is required")
	}

	// Convert account group IDs
	accountGroupIDs := make([]uuid.UUID, 0, len(p.AccountGroupIDs))
	for _, idBytes := range p.AccountGroupIDs {
		id, err := uuid.FromBytes(idBytes)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		accountGroupIDs = append(accountGroupIDs, id)
	}

	return GetTagsReq{
		AccountGroupIDs: accountGroupIDs,
	}, nil
}
