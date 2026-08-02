package model

import (
	"github.com/google/uuid"

	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/finfix/go-server-grpc/proto"
)

type Tag struct {
	ID             uuid.UUID `json:"id" db:"id" minimum:"1"`               // Идентификатор подкатегории
	AccountGroupID uuid.UUID `json:"accountGroupID" db:"account_group_id"` // Идентификатор группы счетов
	Name           string    `json:"name" db:"name"`                       // Название подкатегории
	DatetimeCreate time.Time `json:"datetimeCreate" db:"datetime_create"`  // Дата и время создания
}

// ConvertToProto converts Tag to proto format
func (t Tag) ConvertToProto() (*proto.Tag, error) {
	return &proto.Tag{
		Id:             t.ID[:],
		AccountGroupID: t.AccountGroupID[:],
		Name:           t.Name,
		DatetimeCreate: timestamppb.New(t.DatetimeCreate),
	}, nil
}
