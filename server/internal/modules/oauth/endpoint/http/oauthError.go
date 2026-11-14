package http

import (
	stdErrors "errors"

	pkgErrors "pkg/errors"
	"pkg/log"
	"server/internal/modules/oauth/model"
)

// oauthError извлекает ошибку протокола; любая другая ошибка — server_error с записью в лог
func oauthError(err error) *model.OAuthError {
	var oauthErr *model.OAuthError
	if stdErrors.As(err, &oauthErr) {
		return oauthErr
	}
	log.LogError(pkgErrors.CastError(err))
	return model.NewOAuthError("server_error", "internal error")
}
