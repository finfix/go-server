package model

type SendNotificationReq struct {
	Notification      NotificationSettings
	NotificationToken string
	BundleID          string
}
