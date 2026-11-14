// Package http — HTTP-эндпоинты Authorization Server OAuth 2.1 для MCP-клиентов: метаданные,
// динамическая регистрация, форма входа с выдачей кода, обмен кода и refresh-токена
package http

import (
	"context"

	"server/internal/modules/oauth/model"
	oauthService "server/internal/modules/oauth/service"
)

var _ OAuthService = new(oauthService.OAuthService)

// OAuthService — интерфейс Authorization Server
type OAuthService interface {
	PublicURL() string
	RegisterClient(context.Context, model.RegisterClientReq) (model.RegisterClientRes, error)
	ValidateAuthorize(context.Context, model.AuthorizeReq) (model.Client, error)
	GrantCode(context.Context, model.GrantCodeReq) (string, error)
	ExchangeCode(context.Context, model.ExchangeCodeReq) (model.TokenRes, error)
	Refresh(ctx context.Context, refreshToken string) (model.TokenRes, error)
}

// OAuthEndpoint — HTTP-обработчики OAuth
type OAuthEndpoint struct {
	oauthService OAuthService
}

// NewOAuthEndpoint создает обработчики OAuth
func NewOAuthEndpoint(oauthService OAuthService) *OAuthEndpoint {
	return &OAuthEndpoint{oauthService: oauthService}
}
