package model

import (
	userModel "server/internal/modules/user/model"
	"server/internal/utils/necessary"
)

type RefreshTokensReq struct {
	Token       string                             `json:"token" validate:"required"` // Токен восстановления доступа
	Application userModel.ApplicationInformation   `json:"application"`               // Информация о приложении
	Device      userModel.DeviceInformation        `json:"device"`                    // Информация о девайсе
	Necessary   necessary.NecessaryUserInformation `json:"-"`
}
