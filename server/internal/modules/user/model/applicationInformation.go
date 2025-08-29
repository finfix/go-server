package model

type ApplicationInformation struct {
	BundleID string `json:"bundleID" validate:"required" db:"application_bundle_id"` // Бандл приложения
	Version  string `json:"version" validate:"required" db:"application_version"`    // Версия приложения
	Build    string `json:"build" validate:"required" db:"application_build"`        // Билд приложения
}
