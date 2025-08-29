package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/transaction/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// UpdateTransaction обновление транзакции
func (e *TransactionEndpoint) UpdateTransaction(ctx context.Context, r *proto.UpdateTransactionRequest) (*proto.UpdateTransactionResponse, error) {
	res := &proto.UpdateTransactionResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoUpdateTransactionReq{UpdateTransactionRequest: r}.ConvertToModel()
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
	err = e.transactionService.UpdateTransaction(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}
