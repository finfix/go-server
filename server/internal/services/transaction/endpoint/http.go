package endpoint

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	"pkg/http/chain"
	"server/internal/utils/auth"

	"server/internal/services/transaction/model"
)

var tracer = otel.Tracer("/server/internal/services/transaction/endpoint")

type endpoint struct {
	service transactionService
}

type transactionService interface {
	CreateTransaction(context.Context, model.CreateTransactionReq) (uuid.UUID, error)
	GetTransactions(context.Context, model.GetTransactionsReq) ([]model.Transaction, error)
	UpdateTransaction(context.Context, model.UpdateTransactionReq) error
	DeleteTransaction(context.Context, model.DeleteTransactionReq) error
}

func MountTransactionEndpoints(mux *chi.Mux, service transactionService) {
	mux.Mount("/transaction", newTransactionEndpoint(service))
}

func newTransactionEndpoint(service transactionService) http.Handler {

	s := &endpoint{
		service: service,
	}

	options := []chain.Option{
		chain.Before(auth.DefaultAuthorization),
	}

	router := chi.NewRouter()

	router.Method(http.MethodPost, "/", chain.NewChain(s.createTransaction, options...))
	router.Method(http.MethodPatch, "/", chain.NewChain(s.updateTransaction, options...))
	router.Method(http.MethodDelete, "/", chain.NewChain(s.deleteTransaction, options...))
	router.Method(http.MethodGet, "/", chain.NewChain(s.getTransactions, options...))
	return router
}
