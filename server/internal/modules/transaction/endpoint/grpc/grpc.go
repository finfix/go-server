package grpc

import (
	"context"

	"server/internal/modules/transaction/model"
	transactionService "server/internal/modules/transaction/service"

	proto "github.com/finfix/go-server-grpc/proto"
)

var _ TransactionService = new(transactionService.TransactionService)

type TransactionService interface {
	CreateTransaction(context.Context, model.CreateTransactionReq) (model.CreateTransactionRes, error)
	GetTransactions(context.Context, model.GetTransactionsReq) ([]model.Transaction, error)
	UpdateTransaction(context.Context, model.UpdateTransactionReq) error
	DeleteTransaction(context.Context, model.DeleteTransactionReq) error
}

var _ proto.TransactionEndpointServer = new(TransactionEndpoint)

type TransactionEndpoint struct {
	proto.UnsafeTransactionEndpointServer
	transactionService TransactionService
}

func NewTransactionEndpoint(transactionService TransactionService) *TransactionEndpoint {
	return &TransactionEndpoint{
		transactionService: transactionService,
	}
}
