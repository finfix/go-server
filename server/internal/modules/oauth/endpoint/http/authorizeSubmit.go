package http

import (
	stdErrors "errors"
	"net/http"
	"net/url"

	"server/internal/modules/oauth/model"
)

// AuthorizeSubmit проверяет пароль и редиректит клиенту с кодом; при неверном пароле — форма снова
func (e *OAuthEndpoint) AuthorizeSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeOAuthError(w, model.NewOAuthError("invalid_request", "bad form"))
		return
	}
	req := authorizeReq(r.PostForm)

	code, err := e.oauthService.GrantCode(r.Context(), model.GrantCodeReq{
		Authorize: req,
		Email:     r.PostForm.Get("email"),
		Password:  r.PostForm.Get("password"),
	})
	if err != nil {
		// Неверный пароль — остаемся на форме, остальное — как при GET
		if stdErrors.Is(err, model.ErrInvalidCredentials) {
			client, validateErr := e.oauthService.ValidateAuthorize(r.Context(), req)
			if validateErr != nil {
				e.authorizeError(w, req, validateErr)
				return
			}
			renderForm(w, http.StatusUnauthorized, formData{Client: client.Name, Params: r.PostForm, Error: "Неверная почта или пароль"})
			return
		}
		e.authorizeError(w, req, err)
		return
	}

	// Успех: код и state — клиенту
	redirect, _ := url.Parse(req.RedirectURI)
	query := redirect.Query()
	query.Set("code", code)
	if req.State != "" {
		query.Set("state", req.State)
	}
	redirect.RawQuery = query.Encode()
	http.Redirect(w, r, redirect.String(), http.StatusFound)
}
