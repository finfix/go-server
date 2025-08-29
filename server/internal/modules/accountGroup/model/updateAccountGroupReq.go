package model

import (
	"github.com/google/uuid"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	"github.com/finfix/go-server-grpc/proto"
)

type UpdateAccountGroupReq struct {
	Necessary    necessary.NecessaryUserInformation
	ID           uuid.UUID `json:"id" db:"id"`                       // Идентификатор группы счетов
	Name         *string   `json:"name" db:"name"`                   // Название группы счетов
	Currency     *string   `json:"currency" db:"currency_signatura"` // Валюта группы счетов
	Visible      *bool     `json:"visible" db:"visible"`             // Видимость группы счетов
	SerialNumber *uint32   `json:"serialNumber" db:"serial_number"`  // Порядковый номер группы счетов
}

// ProtoUpdateAccountGroupReq wrapper for proto request
type ProtoUpdateAccountGroupReq struct {
	*proto.UpdateAccountGroupRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoUpdateAccountGroupReq) ConvertToModel() (UpdateAccountGroupReq, error) {
	var res UpdateAccountGroupReq

	if p.UpdateAccountGroupRequest == nil {
		return res, errors.BadRequest.New("UpdateAccountGroupRequest is required")
	}

	// Parse ID from bytes
	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	return UpdateAccountGroupReq{
		ID:           id,
		Name:         p.Name,
		Currency:     p.Currency,
		Visible:      p.Visible,
		SerialNumber: p.SerialNumber,
	}, nil
}
