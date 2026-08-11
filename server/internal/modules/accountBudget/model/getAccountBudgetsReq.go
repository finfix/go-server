package model

import (
	"time"

	"github.com/google/uuid"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	"github.com/finfix/go-server-grpc/proto"
)

// GetAccountBudgetsReq - запрос на получение всех версий бюджета по всем доступным пользователю счетам
type GetAccountBudgetsReq struct {
	Necessary       necessary.NecessaryUserInformation
	AccountGroupIDs []uuid.UUID `json:"accountGroupIDs"` // Идентификаторы групп счетов (пусто - все доступные пользователю группы)
	DateFrom        *time.Time  `json:"dateFrom"`        // Дата, от которой показывать версии бюджета
	DateTo          *time.Time  `json:"dateTo"`          // Дата, до которой показывать версии бюджета
}

// ProtoGetAccountBudgetsReq wrapper for proto request
type ProtoGetAccountBudgetsReq struct {
	*proto.GetAccountBudgetsRequest
}

// ConvertToModel converts proto request to internal model
func (p ProtoGetAccountBudgetsReq) ConvertToModel() (res GetAccountBudgetsReq, err error) {
	if p.GetAccountBudgetsRequest == nil {
		return res, errors.BadRequest.New("GetAccountBudgetsRequest is required")
	}

	var dateFrom *time.Time
	if p.DateFrom != nil {
		t := p.DateFrom.AsTime()
		dateFrom = &t
	}

	var dateTo *time.Time
	if p.DateTo != nil {
		t := p.DateTo.AsTime()
		dateTo = &t
	}

	// Convert account group IDs
	accountGroupIDs := make([]uuid.UUID, 0, len(p.AccountGroupIDs))
	for _, idBytes := range p.AccountGroupIDs {
		id, err := uuid.FromBytes(idBytes)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		accountGroupIDs = append(accountGroupIDs, id)
	}

	return GetAccountBudgetsReq{
		AccountGroupIDs: accountGroupIDs,
		DateFrom:        dateFrom,
		DateTo:          dateTo,
	}, nil
}

// GetAccountBudgetsRes - ответ на получение всех версий бюджета по всем доступным счетам
type GetAccountBudgetsRes struct {
	Budgets []AccountBudget
}

// ConvertToProto преобразует GetAccountBudgetsRes в proto-формат
func (s *GetAccountBudgetsRes) ConvertToProto() *proto.GetAccountBudgetsResponse {
	protoBudgets := make([]*proto.AccountBudget, 0, len(s.Budgets))
	for _, budget := range s.Budgets {
		protoBudgets = append(protoBudgets, budget.ConvertToProto())
	}

	return &proto.GetAccountBudgetsResponse{
		Error:   nil,
		Budgets: protoBudgets,
	}
}
