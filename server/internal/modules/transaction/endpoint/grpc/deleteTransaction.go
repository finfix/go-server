package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/transaction/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// DeleteTransaction удаление транзакции
func (e *TransactionEndpoint) DeleteTransaction(ctx context.Context, r *proto.DeleteTransactionRequest) (*proto.DeleteTransactionResponse, error) {
	res := &proto.DeleteTransactionResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoDeleteTransactionReq{DeleteTransactionRequest: r}.ConvertToModel()
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
	err = e.transactionService.DeleteTransaction(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}
