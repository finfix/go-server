package model

type Version struct {
	Version string `json:"version"` // Версия приложения
	Build   string `json:"build"`   // Номер сборки
}
