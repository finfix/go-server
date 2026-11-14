package mcp

import (
	"net/http"

	sdkAuth "github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/oauthex"
)

// NewResourceMetadataHandler отдает метаданные защищенного ресурса: клиент узнает из них, где
// Authorization Server. У нас это тот же хост
func NewResourceMetadataHandler(publicURL string) http.Handler {
	return sdkAuth.ProtectedResourceMetadataHandler(&oauthex.ProtectedResourceMetadata{ //nolint:exhaustruct
		Resource:               publicURL,
		AuthorizationServers:   []string{publicURL},
		BearerMethodsSupported: []string{"header"},
		ResourceName:           "Coin",
	})
}
