package service

import (
	"net/url"

	"server/internal/modules/oauth/model"
)

// validateRedirectURI проверяет, что адрес абсолютный и либо https, либо http на локальной машине
func validateRedirectURI(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil || !parsed.IsAbs() || parsed.Host == "" {
		return model.NewOAuthError("invalid_redirect_uri", "redirect_uri must be an absolute URL: "+raw)
	}
	switch {
	case parsed.Scheme == "https":
	case parsed.Scheme == "http" && (parsed.Hostname() == "localhost" || parsed.Hostname() == "127.0.0.1" || parsed.Hostname() == "::1"):
	default:
		return model.NewOAuthError("invalid_redirect_uri", "redirect_uri must use https or http on localhost: "+raw)
	}
	return nil
}
