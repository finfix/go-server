package endpoint

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"

	"pkg/http/chain"
	"server/internal/utils/auth"

	"server/internal/modules/settings/model"
	"server/internal/modules/settings/model/applicationType"
)

var tracer = otel.Tracer("/server/internal/modules/settings/endpoint")

type endpoint struct {
	service settingsService
}

type settingsService interface {
	UpdateCurrencies(context.Context, model.UpdateCurrenciesReq) error
	SendNotification(context.Context, model.SendNotificationReq) (model.SendNotificationRes, error)
	GetCurrencies(context.Context) ([]model.Currency, error)
	GetIcons(context.Context) ([]model.Icon, error)
	GetVersion(context.Context, applicationType.Type) (model.Version, error)
}

func MountSettingsEndpoints(mux *chi.Mux, service settingsService) {
	mux.Mount("/settings", newSettingsEndpoint(service))
}

func newSettingsEndpoint(service settingsService) http.Handler {

	e := &endpoint{
		service: service,
	}

	options := []chain.Option{
		chain.Before(auth.DefaultAuthorization),
	}

	r := chi.NewRouter()

	r.Method(http.MethodPost, "/updateCurrencies", chain.NewChain(e.updateCurrencies, options...))
	r.Method(http.MethodPost, "/sendNotification", chain.NewChain(e.sendNotification, options...))
	r.Method(http.MethodGet, "/currencies", chain.NewChain(e.getCurrencies, options...))
	r.Method(http.MethodGet, "/icons", chain.NewChain(e.getIcons, options...))

	// Without authorization
	r.Method(http.MethodGet, "/version/{application}", chain.NewChain(e.getVersion))

	return r
}
