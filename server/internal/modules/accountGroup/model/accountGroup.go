package model

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"
)

type AccountGroup struct {
	ID             uuid.UUID `json:"id" db:"id"`                          // Идентификатор группы счетов
	Name           string    `json:"name" db:"name"`                      // Название группы счетов
	Currency       string    `json:"currency" db:"currency_signatura"`    // Валюта группы счетов
	SerialNumber   uint32    `json:"serialNumber" db:"serial_number"`     // Порядковый номер группы счетов
	Visible        bool      `json:"visible" db:"visible"`                // Видимость группы счетов
	DatetimeCreate time.Time `json:"datetimeCreate" db:"datetime_create"` // Дата и время создания группы счетов
}

func (a *AccountGroup) ConvertToProto() (*proto.AccountGroup, error) {
	return &proto.AccountGroup{
		Id:             a.ID[:],
		Name:           a.Name,
		Currency:       a.Currency,
		Visible:        a.Visible,
		SerialNumber:   a.SerialNumber,
		DatetimeCreate: timestamppb.New(a.DatetimeCreate),
	}, nil
}
