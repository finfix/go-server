package http

import (
	"net/http"
)

// AuthorizeForm проверяет запрос авторизации и показывает форму входа
func (e *OAuthEndpoint) AuthorizeForm(w http.ResponseWriter, r *http.Request) {
	req := authorizeReq(r.URL.Query())
	client, err := e.oauthService.ValidateAuthorize(r.Context(), req)
	if err != nil {
		e.authorizeError(w, req, err)
		return
	}
	renderForm(w, http.StatusOK, formData{Client: client.Name, Params: r.URL.Query(), Error: ""})
}
