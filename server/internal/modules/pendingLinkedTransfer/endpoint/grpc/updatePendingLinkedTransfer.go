package grpc

import (
	"context"

	"pkg/validator"

	"server/internal/modules/pendingLinkedTransfer/model"
	"server/internal/utils/necessary"

	proto "github.com/finfix/go-server-grpc/proto"
)

// UpdatePendingLinkedTransfer меняет статус переноса (Completed/Ignored)
func (e *PendingLinkedTransferEndpoint) UpdatePendingLinkedTransfer(ctx context.Context, r *proto.UpdatePendingLinkedTransferRequest) (*proto.UpdatePendingLinkedTransferResponse, error) {
	res := &proto.UpdatePendingLinkedTransferResponse{}

	req, err := model.ProtoUpdatePendingLinkedTransferReq{UpdatePendingLinkedTransferRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	if err := validator.Validate(req); err != nil {
		return res, err
	}

	if err := req.Validate(); err != nil {
		return res, err
	}

	if err := e.pendingLinkedTransferService.UpdatePendingLinkedTransfer(ctx, req); err != nil {
		return res, err
	}

	return res, nil
}
