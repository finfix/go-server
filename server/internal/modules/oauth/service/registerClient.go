package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"server/internal/modules/oauth/model"
)

// RegisterClient регистрирует клиента динамически. Регистрация открытая: секрета у клиента нет,
// безопасность держится на входе пользователя при авторизации
func (s *OAuthService) RegisterClient(ctx context.Context, req model.RegisterClientReq) (model.RegisterClientRes, error) {
	ctx, span := tracer.Start(ctx, "RegisterClient")
	defer span.End()

	// redirect_uri обязательны и должны быть абсолютными: https либо http на localhost
	if len(req.RedirectURIs) == 0 {
		return model.RegisterClientRes{}, model.NewOAuthError("invalid_redirect_uri", "redirect_uris is required")
	}
	for _, raw := range req.RedirectURIs {
		if err := validateRedirectURI(raw); err != nil {
			return model.RegisterClientRes{}, err
		}
	}
	name := req.ClientName
	if name == "" {
		name = "MCP-клиент"
	}

	// Сохраняем клиента
	id := uuid.New()
	err := s.oauthRepository.CreateClient(ctx, model.CreateClientReq{ID: id, Name: name, RedirectURIs: req.RedirectURIs, CreatedAt: time.Now()})
	if err != nil {
		return model.RegisterClientRes{}, err
	}

	return model.RegisterClientRes{
		ClientID:                id.String(),
		ClientName:              name,
		RedirectURIs:            req.RedirectURIs,
		TokenEndpointAuthMethod: "none",
		GrantTypes:              []string{"authorization_code", "refresh_token"},
		ResponseTypes:           []string{"code"},
	}, nil
}
