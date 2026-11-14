package mcp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	sdkAuth "github.com/modelcontextprotocol/go-sdk/auth"

	"pkg/jwtManager"
	userModel "server/internal/modules/user/model"
	"server/internal/utils/auth"
)

// DeviceIDKey — ключ идентификатора устройства в TokenInfo.Extra
const DeviceIDKey = "deviceID"

// DeviceChecker — проверка, что устройство пользователя еще существует
type DeviceChecker interface {
	GetDevice(ctx context.Context, userID uuid.UUID, deviceID string) (userModel.Device, error)
}

// NewTokenVerifier проверяет Bearer-токен как access-JWT сервера: подпись и срок проверяет
// jwtManager, поэтому Expiration в TokenInfo не заполняется (AllowMissingExpiration). Устройство
// должно существовать: удаление устройства отзывает доступ сразу, не дожидаясь истечения токена
func NewTokenVerifier(deviceChecker DeviceChecker) sdkAuth.TokenVerifier {
	return func(ctx context.Context, token string, _ *http.Request) (*sdkAuth.TokenInfo, error) {
		claims, err := jwtManager.ParseToken[auth.Claims](ctx, token, jwtManager.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", sdkAuth.ErrInvalidToken, err.Error())
		}
		if err = claims.Validate(ctx); err != nil {
			return nil, fmt.Errorf("%w: %s", sdkAuth.ErrInvalidToken, err.Error())
		}

		// Отозванное устройство — токен недействителен
		if _, err = deviceChecker.GetDevice(ctx, claims.UserID, claims.DeviceID); err != nil {
			return nil, fmt.Errorf("%w: device is revoked", sdkAuth.ErrInvalidToken)
		}
		return &sdkAuth.TokenInfo{ //nolint:exhaustruct
			UserID: claims.UserID.String(),
			Extra:  map[string]any{DeviceIDKey: claims.DeviceID},
		}, nil
	}
}
