package grpc

import (
	"context"

	"server/internal/modules/pendingLinkedTransfer/model"
	pendingLinkedTransferService "server/internal/modules/pendingLinkedTransfer/service"

	proto "github.com/finfix/go-server-grpc/proto"
)

var _ PendingLinkedTransferService = new(pendingLinkedTransferService.PendingLinkedTransferService)

type PendingLinkedTransferService interface {
	CreatePendingLinkedTransfer(context.Context, model.CreatePendingLinkedTransferReq) error
	GetPendingLinkedTransfers(context.Context, model.GetPendingLinkedTransfersReq) ([]model.PendingLinkedTransfer, error)
	UpdatePendingLinkedTransfer(context.Context, model.UpdatePendingLinkedTransferReq) error
	DeletePendingLinkedTransfer(context.Context, model.DeletePendingLinkedTransferReq) error
}

var _ proto.PendingLinkedTransferEndpointServer = new(PendingLinkedTransferEndpoint)

type PendingLinkedTransferEndpoint struct {
	proto.UnsafePendingLinkedTransferEndpointServer
	pendingLinkedTransferService PendingLinkedTransferService
}

func NewPendingLinkedTransferEndpoint(pendingLinkedTransferService PendingLinkedTransferService) *PendingLinkedTransferEndpoint {
	return &PendingLinkedTransferEndpoint{
		pendingLinkedTransferService: pendingLinkedTransferService,
	}
}
