package grpc

import (
	"context"

	"pkg/validator"

	"server/internal/modules/pendingLinkedTransfer/model"
	"server/internal/utils/necessary"

	proto "github.com/finfix/go-server-grpc/proto"
)

// CreatePendingLinkedTransfer создаёт требование довнесения транзакции через счёт-мост
func (e *PendingLinkedTransferEndpoint) CreatePendingLinkedTransfer(ctx context.Context, r *proto.CreatePendingLinkedTransferRequest) (*proto.CreatePendingLinkedTransferResponse, error) {
	res := &proto.CreatePendingLinkedTransferResponse{}

	req, err := model.ProtoCreatePendingLinkedTransferReq{CreatePendingLinkedTransferRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	if err := validator.Validate(req); err != nil {
		return res, err
	}

	if err := e.pendingLinkedTransferService.CreatePendingLinkedTransfer(ctx, req); err != nil {
		return res, err
	}

	return res, nil
}
