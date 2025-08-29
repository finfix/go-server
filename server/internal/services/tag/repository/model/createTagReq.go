package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateTagReq struct {
	CreatedByUserID uuid.UUID // Идентификатор пользователя создавшего транзакцию
	AccountGroupID  uuid.UUID // Идентификатор группы счетов
	Name            string    // Название подкатегории
	DatetimeCreate  time.Time // Дата и время создания
}
