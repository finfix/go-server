package model

import (
	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"

	repoModel "server/internal/modules/pendingLinkedTransfer/repository/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type CreatePendingLinkedTransferReq struct {
	Necessary           necessary.NecessaryUserInformation
	ID                  uuid.UUID `json:"id" validate:"required"`                  // Идентификатор переноса
	SourceTransactionID uuid.UUID `json:"sourceTransactionID" validate:"required"` // Транзакция-инициатор
	SourceAccountID     uuid.UUID `json:"sourceAccountID" validate:"required"`     // Счёт-мост со стороны источника
	TargetAccountID     uuid.UUID `json:"targetAccountID" validate:"required"`     // Счёт-мост со стороны получателя
	AccountGroupID      uuid.UUID `json:"accountGroupID" validate:"required"`      // Группа-источник
}

func (s CreatePendingLinkedTransferReq) ConvertToRepoReq() repoModel.CreatePendingLinkedTransferReq {
	return repoModel.CreatePendingLinkedTransferReq{
		ID:                  s.ID,
		SourceTransactionID: s.SourceTransactionID,
		SourceAccountID:     s.SourceAccountID,
		TargetAccountID:     s.TargetAccountID,
		AccountGroupID:      s.AccountGroupID,
	}
}

// ProtoCreatePendingLinkedTransferReq wrapper for proto request
type ProtoCreatePendingLinkedTransferReq struct {
	*proto.CreatePendingLinkedTransferRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoCreatePendingLinkedTransferReq) ConvertToModel() (CreatePendingLinkedTransferReq, error) {
	var res CreatePendingLinkedTransferReq

	if p.CreatePendingLinkedTransferRequest == nil {
		return res, errors.BadRequest.New("CreatePendingLinkedTransferRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	sourceTransactionID, err := uuid.FromBytes(p.SourceTransactionID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	sourceAccountID, err := uuid.FromBytes(p.SourceAccountID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	targetAccountID, err := uuid.FromBytes(p.TargetAccountID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	accountGroupID, err := uuid.FromBytes(p.AccountGroupID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	return CreatePendingLinkedTransferReq{
		ID:                  id,
		SourceTransactionID: sourceTransactionID,
		SourceAccountID:     sourceAccountID,
		TargetAccountID:     targetAccountID,
		AccountGroupID:      accountGroupID,
	}, nil
}
