package model

import (
	"pkg/datetime"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/google/uuid"

	repoModel "server/internal/modules/accountGroup/repository/model"
)

type CreateAccountGroupReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"` // Идентификатор группы счетов

	Name           string        `json:"name" db:"name" validate:"required"`                      // Название группы счетов
	Currency       string        `json:"currency" db:"currency_signatura" validate:"required"`    // Валюта группы счетов
	DatetimeCreate datetime.Time `json:"datetimeCreate" db:"datetime_create" validate:"required"` // Дата и время создания группы счетов
}

func (s CreateAccountGroupReq) ConvertToRepoReq() repoModel.CreateAccountGroupReq {
	return repoModel.CreateAccountGroupReq{
		ID:             s.ID,
		UserID:         s.Necessary.UserID,
		Name:           s.Name,
		Currency:       s.Currency,
		Visible:        true,
		DatetimeCreate: s.DatetimeCreate.Time,
	}
}

// ProtoCreateAccountGroupReq wrapper for proto request
type ProtoCreateAccountGroupReq struct {
	*proto.CreateAccountGroupRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoCreateAccountGroupReq) ConvertToModel() (CreateAccountGroupReq, error) {
	var res CreateAccountGroupReq

	if p.CreateAccountGroupRequest == nil {
		return res, errors.BadRequest.New("CreateAccountGroupRequest is required")
	}

	// Parse ID from bytes
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Convert datetime
	if p.DatetimeCreate == nil {
		return res, errors.BadRequest.New("DatetimeCreate is required")
	}
	datetimeCreate := datetime.Time{Time: p.DatetimeCreate.AsTime()}

	return CreateAccountGroupReq{
		ID:             id,
		Name:           p.Name,
		Currency:       p.Currency,
		DatetimeCreate: datetimeCreate,
	}, nil
}
