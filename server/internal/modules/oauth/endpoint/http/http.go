// Package http — HTTP-эндпоинты Authorization Server OAuth 2.1 для MCP-клиентов: метаданные,
// динамическая регистрация, форма входа с выдачей кода, обмен кода и refresh-токена
package http

import (
	"context"
	"time"

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

// Лимиты защиты от перебора пароля и от засорения базы регистрациями
const (
	loginAttemptsPerEmail = 10        // Попыток входа на одну почту за окно
	loginAttemptsPerIP    = 30        // Попыток входа с одного адреса за окно
	registrationsPerIP    = 20        // Регистраций клиентов с одного адреса за окно
	limiterWindow         = time.Hour // Окно лимитов
)

// OAuthEndpoint — HTTP-обработчики OAuth
type OAuthEndpoint struct {
	oauthService      OAuthService
	loginByEmail      *attemptLimiter // Попытки входа по почте
	loginByIP         *attemptLimiter // Попытки входа по адресу
	registrationsByIP *attemptLimiter // Регистрации клиентов по адресу
}

// NewOAuthEndpoint создает обработчики OAuth
func NewOAuthEndpoint(oauthService OAuthService) *OAuthEndpoint {
	return &OAuthEndpoint{
		oauthService:      oauthService,
		loginByEmail:      newAttemptLimiter(loginAttemptsPerEmail, limiterWindow),
		loginByIP:         newAttemptLimiter(loginAttemptsPerIP, limiterWindow),
		registrationsByIP: newAttemptLimiter(registrationsPerIP, limiterWindow),
	}
}
