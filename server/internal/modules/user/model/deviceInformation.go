package model

import "server/internal/modules/user/model/OS"

type DeviceInformation struct {
	NameOS     OS.OS  `json:"nameOS" validate:"required" db:"device_os_name"`       // Название операционной системы
	VersionOS  string `json:"versionOS" validate:"required" db:"device_os_version"` // Версия операционной системы
	DeviceName string `json:"deviceName" validate:"required" db:"device_name"`      // Название девайса
	ModelName  string `json:"modelName" validate:"required" db:"device_model_name"` // Название модели
	IPAddress  string `json:"-" db:"device_ip_address"`                             // IP-адрес
	UserAgent  string `json:"-" db:"device_user_agent"`                             // UserAgent
}
