package service

import (
	"server/internal/enum/osType"
	userModel "server/internal/modules/user/model"
)

// deviceInformation — устройство OAuth-клиента: без локальной копии данных
func (s *OAuthService) deviceInformation(clientName string) userModel.DeviceInformation {
	return userModel.DeviceInformation{ //nolint:exhaustruct
		NameOS:     osType.Server,
		VersionOS:  s.config.Version,
		DeviceName: clientName,
		ModelName:  "mcp-oauth",
	}
}
