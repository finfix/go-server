package mcp

import (
	"net/http"

	"github.com/google/uuid"
	sdkAuth "github.com/modelcontextprotocol/go-sdk/auth"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"server/internal/utils/contextKeys"
)

// ResourceMetadataPath — путь метаданных защищенного ресурса (RFC 9728)
const ResourceMetadataPath = "/.well-known/oauth-protected-resource"

// NewHandler оборачивает MCP-сервер в HTTP: Bearer-токен проверяется как JWT сервера, без токена
// клиент получает 401 с адресом метаданных ресурса и начинает OAuth. Затем Streamable HTTP без
// сессий — каждый запрос независим, как и остальной сервер
func NewHandler(server *sdk.Server, publicURL string, deviceChecker DeviceChecker) http.Handler {
	transport := sdk.NewStreamableHTTPHandler(
		func(*http.Request) *sdk.Server { return server },
		&sdk.StreamableHTTPOptions{ //nolint:exhaustruct
			Stateless:    true,
			JSONResponse: true,
		},
	)

	// Пользователь и устройство из токена — в контекст, как после gRPC-интерсептора авторизации
	withNecessary := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := sdkAuth.TokenInfoFromContext(r.Context())
		userID, err := uuid.Parse(info.UserID)
		deviceID, _ := info.Extra[DeviceIDKey].(string)
		if err != nil || deviceID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := contextKeys.SetUserID(r.Context(), userID)
		ctx = contextKeys.SetDeviceID(ctx, deviceID)
		transport.ServeHTTP(w, r.WithContext(ctx))
	})

	requireToken := sdkAuth.RequireBearerToken(NewTokenVerifier(deviceChecker), &sdkAuth.RequireBearerTokenOptions{ //nolint:exhaustruct
		ResourceMetadataURL:    publicURL + ResourceMetadataPath,
		AllowMissingExpiration: true,
	})
	return requireToken(withNecessary)
}
