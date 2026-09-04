package model

import "github.com/google/uuid"

type CreatePendingLinkedTransferReq struct {
	ID                  uuid.UUID
	SourceTransactionID uuid.UUID
	SourceAccountID     uuid.UUID
	TargetAccountID     uuid.UUID
	AccountGroupID      uuid.UUID
}
