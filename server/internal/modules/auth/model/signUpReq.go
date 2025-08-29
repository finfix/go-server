package model

import userModel "server/internal/modules/user/model"

type SignUpReq struct {
	Name        string                           `json:"name" validate:"required"`                 // Имя пользователя
	Email       string                           `json:"email" validate:"required" format:"email"` // Электронная почта пользователя
	Password    string                           `json:"password" validate:"required"`             // Пароль пользователя
	Application userModel.ApplicationInformation `json:"application"`                              // Информация о приложении
	Device      userModel.DeviceInformation      `json:"device"`                                   // Информация о девайсе
	DeviceID    string                           `json:"-" validate:"required"`                    // Идентификатор устройства
}
