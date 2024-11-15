package endpoint

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"

	"pkg/http/chain"

	"server/internal/services/accountGroup/model"
)

var tracer = otel.Tracer("/server/internal/services/accountGroup/endpoint")

type endpoint struct {
	service accountGroupService
}

type accountGroupService interface {
	CreateAccountGroup(context.Context, model.CreateAccountGroupReq) (model.CreateAccountGroupRes, error)
	GetAccountGroups(context.Context, model.GetAccountGroupsReq) ([]model.AccountGroup, error)
	UpdateAccountGroup(context.Context, model.UpdateAccountGroupReq) error
	DeleteAccountGroup(context.Context, model.DeleteAccountGroupReq) error
}

func MountAccountGroupEndpoints(mux *chi.Mux, service accountGroupService) {
	mux.Mount("/accountGroup", newAccountGroupEndpoint(service))
}

func newAccountGroupEndpoint(service accountGroupService) http.Handler {

	e := &endpoint{
		service: service,
	}

	options := []chain.Option{
		chain.Before(chain.DefaultAuthorization),
	}

	r := chi.NewRouter()

	r.Method(http.MethodPost, "/", chain.NewChain(e.createAccountGroup, options...))
	r.Method(http.MethodGet, "/", chain.NewChain(e.getAccountGroups, options...))
	r.Method(http.MethodPatch, "/", chain.NewChain(e.updateAccountGroup, options...))
	r.Method(http.MethodDelete, "/", chain.NewChain(e.deleteAccountGroup, options...))

	return r
}
