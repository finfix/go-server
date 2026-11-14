package service

import (
	authModel "server/internal/modules/auth/model"
	"server/internal/modules/oauth/model"
)

// tokenRes собирает ответ эндпоинта токенов
func (s *OAuthService) tokenRes(tokens authModel.Tokens) model.TokenRes {
	return model.TokenRes{
		AccessToken:  tokens.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.config.AccessTokenTTL.Seconds()),
		RefreshToken: tokens.RefreshToken,
	}
}
