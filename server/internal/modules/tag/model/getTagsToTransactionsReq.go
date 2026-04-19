package model

import (
	"github.com/google/uuid"

	"server/internal/utils/necessary"
	"server/internal/utils/errors"
	"github.com/finfix/go-server-grpc/proto"
)

type GetTagsToTransactionsReq struct {
	Necessary       necessary.NecessaryUserInformation
	AccountGroupIDs []uuid.UUID `json:"-" schema:"-" minimum:"1"` // Идентификаторы групп счетов
	TransactionIDs  []uuid.UUID `json:"-" schema:"-" minimum:"1"` // Идентификаторы транзакций
}

// ProtoGetTagsToTransactionsReq wrapper for proto request
type ProtoGetTagsToTransactionsReq struct {
	*proto.GetTagsToTransactionsRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoGetTagsToTransactionsReq) ConvertToModel() (GetTagsToTransactionsReq, error) {
	if p.GetTagsToTransactionsRequest == nil {
		return GetTagsToTransactionsReq{}, errors.BadRequest.New("GetTagsToTransactionsRequest is required")
	}

	// Convert account group IDs
	accountGroupIDs := make([]uuid.UUID, 0, len(p.AccountGroupIDs))
	for _, idBytes := range p.AccountGroupIDs {
		id, err := uuid.FromBytes(idBytes)
		if err != nil {
			return GetTagsToTransactionsReq{}, errors.BadRequest.Wrap(err)
		}
		accountGroupIDs = append(accountGroupIDs, id)
	}

	// Convert transaction IDs
	transactionIDs := make([]uuid.UUID, 0, len(p.TransactionIDs))
	for _, idBytes := range p.TransactionIDs {
		id, err := uuid.FromBytes(idBytes)
		if err != nil {
			return GetTagsToTransactionsReq{}, errors.BadRequest.Wrap(err)
		}
		transactionIDs = append(transactionIDs, id)
	}

	return GetTagsToTransactionsReq{
		AccountGroupIDs: accountGroupIDs,
		TransactionIDs:  transactionIDs,
	}, nil
}
