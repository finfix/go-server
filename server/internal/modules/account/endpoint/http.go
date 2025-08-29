package endpoint

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"

	"pkg/http/chain"
	"server/internal/modules/account/model"
	"server/internal/utils/auth"
)

var tracer = otel.Tracer("/server/internal/modules/account/endpoint")

type endpoint struct {
	service accountService
}

type accountService interface {
	CreateAccount(context.Context, model.CreateAccountReq) (model.CreateAccountRes, error)
	GetAccounts(context.Context, model.GetAccountsReq) ([]model.Account, error)
	UpdateAccount(context.Context, model.UpdateAccountReq) (model.UpdateAccountRes, error)
	DeleteAccount(context.Context, model.DeleteAccountReq) error
}

// MountAccountEndpoints mounts account endpoints to the router
func MountAccountEndpoints(mux *chi.Mux, service accountService) {
	mux.Mount("/account", newAccountEndpoint(service))
}

func newAccountEndpoint(service accountService) http.Handler {

	e := &endpoint{
		service: service,
	}

	options := []chain.Option{
		chain.Before(auth.DefaultAuthorization),
	}

	r := chi.NewRouter()

	r.Method(http.MethodPost, "/", chain.NewChain(e.createAccount, options...))
	r.Method(http.MethodGet, "/", chain.NewChain(e.getAccounts, options...))
	r.Method(http.MethodPatch, "/", chain.NewChain(e.updateAccount, options...))
	r.Method(http.MethodDelete, "/", chain.NewChain(e.deleteAccount, options...))

	return r
}
