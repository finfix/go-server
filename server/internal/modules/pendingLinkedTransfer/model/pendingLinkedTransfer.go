package model

import (
	"time"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"

	"server/internal/enum/pendingLinkedTransferStatus"
)

// PendingLinkedTransfer - требование довнесения транзакции через счёт-мост
type PendingLinkedTransfer struct {
	ID                  uuid.UUID                                                `json:"id" db:"id"`
	Status              pendingLinkedTransferStatus.PendingLinkedTransferStatus `json:"status" db:"status"`
	SourceTransactionID uuid.UUID                                                `json:"sourceTransactionID" db:"source_transaction_id"`
	SourceAccountID     uuid.UUID                                                `json:"sourceAccountID" db:"source_account_id"`
	TargetAccountID     uuid.UUID                                                `json:"targetAccountID" db:"target_account_id"`
	AccountGroupID      uuid.UUID                                                `json:"accountGroupID" db:"account_group_id"`
	DatetimeCreate      time.Time                                                `json:"datetimeCreate" db:"datetime_create"`
}

// ConvertToProto converts internal PendingLinkedTransfer to proto PendingLinkedTransfer
func (p PendingLinkedTransfer) ConvertToProto() (*proto.PendingLinkedTransfer, error) {
	status, err := p.Status.ConvertToProto()
	if err != nil {
		return nil, err
	}

	return &proto.PendingLinkedTransfer{
		Id:                  p.ID[:],
		Status:              status,
		SourceTransactionID: p.SourceTransactionID[:],
		SourceAccountID:     p.SourceAccountID[:],
		TargetAccountID:     p.TargetAccountID[:],
		AccountGroupID:      p.AccountGroupID[:],
	}, nil
}
