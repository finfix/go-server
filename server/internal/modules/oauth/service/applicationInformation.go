package service

import (
	userModel "server/internal/modules/user/model"
)

// applicationInformation — приложение OAuth-клиента для записи устройства
func (s *OAuthService) applicationInformation() userModel.ApplicationInformation {
	return userModel.ApplicationInformation{BundleID: "coin.mcp.oauth", Version: s.config.Version, Build: s.config.Version}
}
