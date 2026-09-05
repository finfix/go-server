package model

import (
	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"

	"server/internal/enum/pendingLinkedTransferStatus"
	repoModel "server/internal/modules/pendingLinkedTransfer/repository/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type GetPendingLinkedTransfersReq struct {
	Necessary        necessary.NecessaryUserInformation
	IDs              []uuid.UUID                                              `json:"-"`
	AccountGroupIDs  []uuid.UUID                                              `json:"accountGroupIDs"`
	TargetAccountIDs []uuid.UUID                                              `json:"targetAccountIDs"`
	Status           *pendingLinkedTransferStatus.PendingLinkedTransferStatus `json:"status"`
}

func (s GetPendingLinkedTransfersReq) ConvertToRepoReq() repoModel.GetPendingLinkedTransfersReq {
	return repoModel.GetPendingLinkedTransfersReq{
		IDs:              s.IDs,
		AccountGroupIDs:  s.AccountGroupIDs,
		TargetAccountIDs: s.TargetAccountIDs,
		Status:           s.Status,
	}
}

// ProtoGetPendingLinkedTransfersReq wrapper for proto request
type ProtoGetPendingLinkedTransfersReq struct {
	*proto.GetPendingLinkedTransfersRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoGetPendingLinkedTransfersReq) ConvertToModel() (GetPendingLinkedTransfersReq, error) {
	var res GetPendingLinkedTransfersReq

	if p.GetPendingLinkedTransfersRequest == nil {
		return res, errors.BadRequest.New("GetPendingLinkedTransfersRequest is required")
	}

	accountGroupIDs := make([]uuid.UUID, 0, len(p.AccountGroupIDs))
	for _, idBytes := range p.AccountGroupIDs {
		id, err := uuid.FromBytes(idBytes)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		accountGroupIDs = append(accountGroupIDs, id)
	}

	targetAccountIDs := make([]uuid.UUID, 0, len(p.TargetAccountIDs))
	for _, idBytes := range p.TargetAccountIDs {
		id, err := uuid.FromBytes(idBytes)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		targetAccountIDs = append(targetAccountIDs, id)
	}

	var status *pendingLinkedTransferStatus.PendingLinkedTransferStatus
	if p.Status != nil {
		_status, err := pendingLinkedTransferStatus.ProtoPendingLinkedTransferStatus{PendingLinkedTransferStatus: *p.Status}.ConvertToModel()
		if err != nil {
			return res, err
		}
		status = &_status
	}

	return GetPendingLinkedTransfersReq{
		AccountGroupIDs:  accountGroupIDs,
		TargetAccountIDs: targetAccountIDs,
		Status:           status,
	}, nil
}
