package http

import (
	"net/http"
)

// Register добавляет маршруты OAuth в mux
func (e *OAuthEndpoint) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /.well-known/oauth-authorization-server", e.Metadata)
	mux.HandleFunc("GET /.well-known/oauth-authorization-server/", e.Metadata)
	mux.HandleFunc("POST /oauth/register", e.RegisterClient)
	mux.HandleFunc("GET /oauth/authorize", e.AuthorizeForm)
	mux.HandleFunc("POST /oauth/authorize", e.AuthorizeSubmit)
	mux.HandleFunc("POST /oauth/token", e.Token)
}
