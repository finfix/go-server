package model

// NewOAuthError создает ошибку протокола
func NewOAuthError(code, description string) *OAuthError {
	return &OAuthError{Code: code, Description: description}
}
