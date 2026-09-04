package model

import (
	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"

	"server/internal/enum/pendingLinkedTransferStatus"
	repoModel "server/internal/modules/pendingLinkedTransfer/repository/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type UpdatePendingLinkedTransferReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID                                                 `json:"id" validate:"required"` // Идентификатор переноса
	Status    *pendingLinkedTransferStatus.PendingLinkedTransferStatus `json:"status"`                  // Новый статус (Completed/Ignored) — patch, отсутствие значит "не менять"
}

func (s UpdatePendingLinkedTransferReq) Validate() error {
	if s.Status == nil {
		return nil
	}
	switch *s.Status {
	case pendingLinkedTransferStatus.Completed, pendingLinkedTransferStatus.Ignored:
	default:
		return errors.BadRequest.New("Status must be Completed or Ignored").
			WithParams("status", *s.Status)
	}
	return nil
}

func (s UpdatePendingLinkedTransferReq) ConvertToRepoReq() repoModel.UpdatePendingLinkedTransferReq {
	return repoModel.UpdatePendingLinkedTransferReq{
		Status: s.Status,
	}
}

// ProtoUpdatePendingLinkedTransferReq wrapper for proto request
type ProtoUpdatePendingLinkedTransferReq struct {
	*proto.UpdatePendingLinkedTransferRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoUpdatePendingLinkedTransferReq) ConvertToModel() (UpdatePendingLinkedTransferReq, error) {
	var res UpdatePendingLinkedTransferReq

	if p.UpdatePendingLinkedTransferRequest == nil {
		return res, errors.BadRequest.New("UpdatePendingLinkedTransferRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	var status *pendingLinkedTransferStatus.PendingLinkedTransferStatus
	if p.Status != nil {
		_status, err := pendingLinkedTransferStatus.ProtoPendingLinkedTransferStatus{PendingLinkedTransferStatus: *p.Status}.ConvertToModel()
		if err != nil {
			return res, err
		}
		status = &_status
	}

	return UpdatePendingLinkedTransferReq{
		ID:     id,
		Status: status,
	}, nil
}
