package model

import (
	"github.com/google/uuid"

	"server/internal/utils/necessary"
	"server/internal/utils/errors"
	"github.com/finfix/go-server-grpc/proto"
)

type GetAccountGroupsReq struct {
	Necessary       necessary.NecessaryUserInformation
	AccountGroupIDs []uuid.UUID `json:"accountGroupIDs" schema:"accountGroupIDs" minimum:"1"` // Идентификаторы групп счетов
}

// ProtoGetAccountGroupsReq wrapper for proto request
type ProtoGetAccountGroupsReq struct {
	*proto.GetAccountGroupsRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoGetAccountGroupsReq) ConvertToModel() (GetAccountGroupsReq, error) {
	var res GetAccountGroupsReq

	if p.GetAccountGroupsRequest == nil {
		return res, errors.BadRequest.New("GetAccountGroupsRequest is required")
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

	return GetAccountGroupsReq{
		AccountGroupIDs: accountGroupIDs,
	}, nil
}
