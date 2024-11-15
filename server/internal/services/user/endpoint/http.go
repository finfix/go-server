package endpoint

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"

	"pkg/http/chain"

	"server/internal/services/user/model"
)

var tracer = otel.Tracer("/server/internal/services/user/endpoint")

type endpoint struct {
	service userService
}

type userService interface {
	GetUsers(context.Context, model.GetUsersReq) ([]model.User, error)
	UpdateUser(context.Context, model.UpdateUserReq) error
}

func MountUserEndpoints(mux *chi.Mux, service userService) {
	mux.Mount("/user", newUserEndpoint(service))
}

func newUserEndpoint(service userService) http.Handler {

	e := &endpoint{
		service: service,
	}

	options := []chain.Option{
		chain.Before(chain.DefaultAuthorization),
	}

	r := chi.NewRouter()

	r.Method(http.MethodGet, "/", chain.NewChain(e.getUser, options...))
	r.Method(http.MethodPatch, "/", chain.NewChain(e.updateUser, options...))
	return r
}
