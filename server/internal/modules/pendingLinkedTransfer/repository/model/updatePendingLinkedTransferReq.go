package model

import "server/internal/enum/pendingLinkedTransferStatus"

type UpdatePendingLinkedTransferReq struct {
	Status *pendingLinkedTransferStatus.PendingLinkedTransferStatus
}
