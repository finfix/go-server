package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/transaction/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// CreateTransaction создание транзакции
func (e *TransactionEndpoint) CreateTransaction(ctx context.Context, r *proto.CreateTransactionRequest) (*proto.CreateTransactionResponse, error) {
	res := &proto.CreateTransactionResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoCreateTransactionReq{CreateTransactionRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Parse necessary information from context
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return res, err
	}

	// Call service method
	createRes, err := e.transactionService.CreateTransaction(ctx, req)
	if err != nil {
		return res, err
	}

	// Convert response to proto and return
	return createRes.ConvertToProto(), nil
}
