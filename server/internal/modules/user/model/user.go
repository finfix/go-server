package model

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/finfix/go-server-grpc/proto"

	"pkg/datetime"
)

type User struct {
	ID              uuid.UUID     `db:"id" json:"id"`                            // Идентификатор пользователя
	Name            string        `db:"name" json:"name"`                        // Имя пользователя
	Email           string        `db:"email" json:"email"`                      // Электронная почта
	PasswordHash    []byte        `db:"password_hash" json:"-"`                  // Хэш пароля
	PasswordSalt    []byte        `db:"password_salt" json:"-"`                  // Соль пароля
	TimeCreate      datetime.Time `db:"time_create" json:"timeCreate"`           // Дата и время создания аккаунта
	DefaultCurrency string        `db:"default_currency" json:"defaultCurrency"` // Валюта по умолчанию
	IsAdmin         bool          `db:"is_admin" json:"-"`                       // Является ли пользователь администратором системы
}

// ConvertToProto converts User to proto format
func (u User) ConvertToProto() (*proto.User, error) {
	return &proto.User{
		Id:              u.ID[:],
		Name:            u.Name,
		Email:           u.Email,
		DefaultCurrency: u.DefaultCurrency,
		TimeCreate:      timestamppb.New(u.TimeCreate.Time),
	}, nil
}
