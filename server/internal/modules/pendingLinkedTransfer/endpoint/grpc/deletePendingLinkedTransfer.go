package grpc

import (
	"context"

	"pkg/validator"

	"server/internal/modules/pendingLinkedTransfer/model"
	"server/internal/utils/necessary"

	proto "github.com/finfix/go-server-grpc/proto"
)

// DeletePendingLinkedTransfer удаляет перенос
func (e *PendingLinkedTransferEndpoint) DeletePendingLinkedTransfer(ctx context.Context, r *proto.DeletePendingLinkedTransferRequest) (*proto.DeletePendingLinkedTransferResponse, error) {
	res := &proto.DeletePendingLinkedTransferResponse{}

	req, err := model.ProtoDeletePendingLinkedTransferReq{DeletePendingLinkedTransferRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	if err := validator.Validate(req); err != nil {
		return res, err
	}

	if err := e.pendingLinkedTransferService.DeletePendingLinkedTransfer(ctx, req); err != nil {
		return res, err
	}

	return res, nil
}
