package endpoint

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"

	"pkg/http/chain"
	"server/internal/utils/auth"

	"server/internal/services/tag/model"
)

var tracer = otel.Tracer("/server/internal/services/tag/endpoint")

type endpoint struct {
	service tagService
}

type tagService interface {
	CreateTag(context.Context, model.CreateTagReq) (uint32, error)
	GetTags(context.Context, model.GetTagsReq) ([]model.Tag, error)
	UpdateTag(context.Context, model.UpdateTagReq) error
	DeleteTag(context.Context, model.DeleteTagReq) error

	GetTagsToTransactions(ctx context.Context, req model.GetTagsToTransactionsReq) ([]model.TagToTransaction, error)
}

func MountTagEndpoints(mux *chi.Mux, service tagService) {
	mux.Mount("/tag", newTagEndpoint(service))
}

func newTagEndpoint(service tagService) http.Handler {

	s := &endpoint{
		service: service,
	}

	options := []chain.Option{
		chain.Before(auth.DefaultAuthorization),
	}

	router := chi.NewRouter()

	router.Method(http.MethodPost, "/", chain.NewChain(s.createTag, options...))
	router.Method(http.MethodGet, "/", chain.NewChain(s.getTags, options...))
	router.Method(http.MethodPatch, "/", chain.NewChain(s.updateTag, options...))
	router.Method(http.MethodDelete, "/", chain.NewChain(s.deleteTag, options...))

	router.Method(http.MethodGet, "/to_transactions", chain.NewChain(s.getTagsToTransaction, options...))
	return router
}
