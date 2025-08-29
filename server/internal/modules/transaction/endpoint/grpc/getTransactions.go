package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/transaction/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// GetTransactions получение всех транзакций
func (e *TransactionEndpoint) GetTransactions(ctx context.Context, r *proto.GetTransactionsRequest) (*proto.GetTransactionsResponse, error) {
	res := &proto.GetTransactionsResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoGetTransactionsReq{GetTransactionsRequest: r}.ConvertToModel()
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
	transactions, err := e.transactionService.GetTransactions(ctx, req)
	if err != nil {
		return res, err
	}

	// Convert transactions to proto format
	protoTransactions := make([]*proto.Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		protoTransaction, err := transaction.ConvertToProto()
		if err != nil {
			return res, err
		}
		protoTransactions = append(protoTransactions, protoTransaction)
	}

	res.Transactions = protoTransactions
	return res, nil
}
