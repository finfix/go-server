package pendingLinkedTransferStatus

import (
	"context"
	"pkg/maps"

	"server/internal/utils/errors"

	"github.com/finfix/go-server-grpc/proto"
)

type PendingLinkedTransferStatus string

// enums:"pending,completed,ignored"
const (
	Pending   PendingLinkedTransferStatus = "pending"
	Completed PendingLinkedTransferStatus = "completed"
	Ignored   PendingLinkedTransferStatus = "ignored"
)

func (t *PendingLinkedTransferStatus) Validate(ctx context.Context) error {
	if t == nil {
		return nil
	}
	switch *t {
	case Pending, Completed, Ignored:
	default:
		return errors.BadRequest.New("Unknown pending linked transfer status").
			SkipThisCall().
			WithContextParams(ctx).
			WithParams("status", *t).
			WithCustomHumanText("Неизвестный статус переноса")
	}
	return nil
}

// mappingProtoToModel содержит соответствие между значениями proto.PendingLinkedTransferStatus и PendingLinkedTransferStatus
var mappingProtoToModel = map[proto.PendingLinkedTransferStatus]PendingLinkedTransferStatus{
	proto.PendingLinkedTransferStatus_Pending:   Pending,
	proto.PendingLinkedTransferStatus_Completed: Completed,
	proto.PendingLinkedTransferStatus_Ignored:   Ignored,
}

// ConvertToProto преобразует PendingLinkedTransferStatus в proto.PendingLinkedTransferStatus
func (b PendingLinkedTransferStatus) ConvertToProto() (status proto.PendingLinkedTransferStatus, err error) {
	mappingModelToProto, err := maps.Revert(mappingProtoToModel)
	if err != nil {
		return 0, err
	}

	protoStatus, ok := mappingModelToProto[b]
	if !ok {
		return protoStatus, errors.BadRequest.New("PendingLinkedTransferStatus undefined")
	}

	return protoStatus, nil
}

type ProtoPendingLinkedTransferStatus struct {
	proto.PendingLinkedTransferStatus
}

// ConvertToModel преобразует ProtoPendingLinkedTransferStatus в PendingLinkedTransferStatus
func (p ProtoPendingLinkedTransferStatus) ConvertToModel() (status PendingLinkedTransferStatus, err error) {
	status, ok := mappingProtoToModel[p.PendingLinkedTransferStatus]
	if !ok {
		return status, errors.BadRequest.New("PendingLinkedTransferStatus undefined")
	}

	return status, nil
}
