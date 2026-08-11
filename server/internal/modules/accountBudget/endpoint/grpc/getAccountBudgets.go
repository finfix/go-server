package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/accountBudget/model"
	"server/internal/utils/necessary"

	proto "github.com/finfix/go-server-grpc/proto"
)

// GetAccountBudgets получение всех версий бюджета по всем счетам доступных пользователю групп счетов
func (e *AccountBudgetEndpoint) GetAccountBudgets(ctx context.Context, r *proto.GetAccountBudgetsRequest) (*proto.GetAccountBudgetsResponse, error) {
	res := new(proto.GetAccountBudgetsResponse)

	// Convert proto request to internal model
	req, err := model.ProtoGetAccountBudgetsReq{GetAccountBudgetsRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Parse necessary information from context
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return res, err
	}

	// Call service method
	budgets, err := e.accountBudgetService.GetAccountBudgets(ctx, req)
	if err != nil {
		return res, err
	}

	_res := model.GetAccountBudgetsRes{Budgets: budgets}

	return _res.ConvertToProto(), nil
}
