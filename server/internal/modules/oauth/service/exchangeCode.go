package service

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/google/uuid"

	authModel "server/internal/modules/auth/model"
	"server/internal/modules/oauth/model"
)

// ExchangeCode обменивает код на токены сервера: проверяет клиента, redirect_uri, PKCE и resource,
// создает для выдачи отдельное устройство без локальной копии данных
func (s *OAuthService) ExchangeCode(ctx context.Context, req model.ExchangeCodeReq) (model.TokenRes, error) {
	ctx, span := tracer.Start(ctx, "ExchangeCode")
	defer span.End()

	// Код одноразовый: ConsumeCode помечает его использованным атомарно
	code, found, err := s.oauthRepository.ConsumeCode(ctx, req.Code, time.Now())
	if err != nil {
		return model.TokenRes{}, err
	}
	if !found {
		return model.TokenRes{}, model.NewOAuthError("invalid_grant", "code is unknown, expired or already used")
	}

	// Код должен обменивать тот же клиент на тот же адрес
	clientID, err := uuid.Parse(req.ClientID)
	if err != nil || clientID != code.ClientID {
		return model.TokenRes{}, model.NewOAuthError("invalid_grant", "client_id does not match the code")
	}
	if req.RedirectURI != code.RedirectURI {
		return model.TokenRes{}, model.NewOAuthError("invalid_grant", "redirect_uri does not match the code")
	}
	if req.Resource != "" && !sameResource(req.Resource, code.Resource) {
		return model.TokenRes{}, model.NewOAuthError("invalid_target", "resource does not match the code")
	}

	// PKCE: SHA256(code_verifier) должен совпасть с code_challenge
	sum := sha256.Sum256([]byte(req.CodeVerifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])
	if req.CodeVerifier == "" || subtle.ConstantTimeCompare([]byte(challenge), []byte(code.CodeChallenge)) != 1 {
		return model.TokenRes{}, model.NewOAuthError("invalid_grant", "code_verifier does not match code_challenge")
	}

	// Имя клиента — в название устройства, чтобы его было видно в списке сессий
	client, _, err := s.oauthRepository.GetClient(ctx, code.ClientID)
	if err != nil {
		return model.TokenRes{}, err
	}

	// Одно устройство на пару клиент–пользователь: повторный вход клиента обновляет его, а не
	// плодит строки; отзыв — удаление устройства, как у телефона
	deviceSum := sha256.Sum256([]byte(code.ClientID.String() + ":" + code.UserID.String()))
	tokens, err := s.authService.IssueTokens(ctx, authModel.IssueTokensReq{
		UserID:      code.UserID,
		DeviceID:    "oauth-" + hex.EncodeToString(deviceSum[:8]),
		Device:      s.deviceInformation(client.Name),
		Application: s.applicationInformation(),
	})
	if err != nil {
		return model.TokenRes{}, err
	}

	return s.tokenRes(tokens), nil
}
