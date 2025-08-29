package model

type Notification struct {
	Title      string `json:"title"`      // Заголовок уведомления
	Subtitle   string `json:"subtitle"`   // Подзаголовок уведомления
	Message    string `json:"message"`    // Сообщение уведомления
	BadgeCount uint8  `json:"badgeCount"` // Индикатор какое значение показывать в бадже
}
