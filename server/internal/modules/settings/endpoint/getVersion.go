package endpoint

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "server/internal/modules/settings/model" //nolint:golint
	"server/internal/modules/settings/model/applicationType"
)

// @Summary Получение текущей версии сервера
// @Tags settings
// @Produce json
// @Success 200 {object} model.Version
// @Router /settings/version/ [get]
func (s *endpoint) getVersion(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "getVersion")
	defer span.End()

	appType := applicationType.Type(chi.URLParam(r, "application"))
	if err := appType.Validate(ctx); err != nil {
		return nil, err
	}
	return s.service.GetVersion(ctx, appType)
}
