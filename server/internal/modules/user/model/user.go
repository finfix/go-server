package model

import (
	"github.com/google/uuid"

	"pkg/datetime"
)

type User struct {
	ID              uuid.UUID     `db:"id" json:"id"`                                      // Идентификатор пользователя
	Name            string        `db:"name" json:"name"`                                  // Имя пользователя
	Email           string        `db:"email" json:"email"`                                // Электронная почта
	PasswordHash    []byte        `db:"password_hash" json:"-"`                            // Хэш пароля
	PasswordSalt    []byte        `db:"password_salt" json:"-"`                            // Соль пароля
	TimeCreate      datetime.Time `db:"time_create" json:"timeCreate"`                     // Дата и время создания аккаунта
	DefaultCurrency string        `db:"default_currency_signatura" json:"defaultCurrency"` // Валюта по умолчанию
	IsAdmin         bool          `db:"is_admin" json:"-"`                                 // Является ли пользователь администратором системы
}
