package model

import (
	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"
)

type Icon struct {
	ID   uuid.UUID `json:"id" db:"id"`     // ID иконки
	Name string    `json:"name" db:"name"` // Название иконки
	Url  string    `json:"url" db:"img"`   // URL иконки
}

func (i *Icon) ConvertToProto() (*proto.Icon, error) {
	return &proto.Icon{
		Id:   i.ID[:],
		Name: i.Name,
		Url:  i.Url,
	}, nil
}
