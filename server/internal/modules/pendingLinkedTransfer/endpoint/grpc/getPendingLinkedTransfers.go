package grpc

import (
	"context"

	"server/internal/modules/pendingLinkedTransfer/model"
	"server/internal/utils/necessary"

	proto "github.com/finfix/go-server-grpc/proto"
)

// GetPendingLinkedTransfers получение переносов по фильтрам
func (e *PendingLinkedTransferEndpoint) GetPendingLinkedTransfers(ctx context.Context, r *proto.GetPendingLinkedTransfersRequest) (*proto.GetPendingLinkedTransfersResponse, error) {
	res := &proto.GetPendingLinkedTransfersResponse{}

	req, err := model.ProtoGetPendingLinkedTransfersReq{GetPendingLinkedTransfersRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	transfers, err := e.pendingLinkedTransferService.GetPendingLinkedTransfers(ctx, req)
	if err != nil {
		return res, err
	}

	for _, transfer := range transfers {
		protoTransfer, err := transfer.ConvertToProto()
		if err != nil {
			return res, err
		}
		res.PendingLinkedTransfers = append(res.PendingLinkedTransfers, protoTransfer)
	}

	return res, nil
}
