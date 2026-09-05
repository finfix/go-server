package model

import (
	"github.com/google/uuid"

	"server/internal/enum/pendingLinkedTransferStatus"
)

type GetPendingLinkedTransfersReq struct {
	IDs              []uuid.UUID
	AccountGroupIDs  []uuid.UUID
	TargetAccountIDs []uuid.UUID
	Status           *pendingLinkedTransferStatus.PendingLinkedTransferStatus
}
