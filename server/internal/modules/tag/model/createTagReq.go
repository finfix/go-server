package model

import (
	"github.com/google/uuid"

	"pkg/datetime"
	"server/internal/utils/necessary"
	"server/internal/utils/errors"
	"github.com/finfix/go-server-grpc/proto"

	repoModel "server/internal/modules/tag/repository/model"
)

type CreateTagReq struct {
	Necessary      necessary.NecessaryUserInformation
	ID             uuid.UUID     `json:"id" validate:"required"`             // Идентификатор тега
	AccountGroupID uuid.UUID     `json:"accountGroupID" validate:"required"` // Идентификатор группы счетов
	Name           string        `json:"name" validate:"required"`           // Название подкатегории
	DatetimeCreate datetime.Time `json:"datetimeCreate" validate:"required"` // Дата создания подкатегории
}

func (s CreateTagReq) ConvertToRepoReq() repoModel.CreateTagReq {
	return repoModel.CreateTagReq{
		ID:              s.ID,
		Name:            s.Name,
		AccountGroupID:  s.AccountGroupID,
		CreatedByUserID: s.Necessary.UserID,
		DatetimeCreate:  s.DatetimeCreate.Time,
	}
}

// ProtoCreateTagReq wrapper for proto request
type ProtoCreateTagReq struct {
	*proto.CreateTagRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoCreateTagReq) ConvertToModel() (CreateTagReq, error) {
	var res CreateTagReq

	if p.CreateTagRequest == nil {
		return res, errors.BadRequest.New("CreateTagRequest is required")
	}

	// Parse ID from bytes
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Parse AccountGroupID
	accountGroupID, err := uuid.FromBytes(p.AccountGroupID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Convert datetime
	if p.DatetimeCreate == nil {
		return res, errors.BadRequest.New("DatetimeCreate is required")
	}
	datetimeCreate := datetime.Time{Time: p.DatetimeCreate.AsTime()}

	return CreateTagReq{
		ID:             id,
		AccountGroupID: accountGroupID,
		Name:           p.Name,
		DatetimeCreate: datetimeCreate,
	}, nil
}
