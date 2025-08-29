package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateAccountGroupReq struct {
	ID             uuid.UUID // Идентификатор группы счетов
	Name           string    // Название группы счетов
	Currency       string    // Валюта группы счетов
	Visible        bool      // Видимость группы счетов
	DatetimeCreate time.Time // Дата и время создания группы счетов
	UserID         uuid.UUID // Каким пользователем создан объект
}
