package http

import (
	"net/http"
	"net/url"

	"server/internal/modules/oauth/model"
)

// authorizeError: ошибку клиента, redirect_uri или сервера показываем пользователю, остальные —
// редиректом клиенту (RFC 6749 §4.1.2.1). При ошибке сервера redirect_uri мог остаться
// непроверенным, поэтому редиректить на него нельзя
func (e *OAuthEndpoint) authorizeError(w http.ResponseWriter, req model.AuthorizeReq, err error) {
	oauthErr := oauthError(err)
	if oauthErr.Code == "invalid_client" || oauthErr.Code == "server_error" || oauthErr.Code == "invalid_request" && req.RedirectURI == "" ||
		oauthErr.Description == "redirect_uri is not registered for this client" {
		renderError(w, http.StatusBadRequest, oauthErr.Description)
		return
	}
	redirect, parseErr := url.Parse(req.RedirectURI)
	if parseErr != nil {
		renderError(w, http.StatusBadRequest, oauthErr.Description)
		return
	}
	query := redirect.Query()
	query.Set("error", oauthErr.Code)
	query.Set("error_description", oauthErr.Description)
	if req.State != "" {
		query.Set("state", req.State)
	}
	redirect.RawQuery = query.Encode()
	w.Header().Set("Location", redirect.String())
	w.WriteHeader(http.StatusFound)
}
