package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/accountBudget/model"
	"server/internal/utils/necessary"

	proto "github.com/finfix/go-server-grpc/proto"
)

// CreateAccountBudget создание новой версии бюджета счета
func (e *AccountBudgetEndpoint) CreateAccountBudget(ctx context.Context, r *proto.CreateAccountBudgetRequest) (*proto.CreateAccountBudgetResponse, error) {
	res := new(proto.CreateAccountBudgetResponse)

	// Convert proto request to internal model
	req, err := model.ProtoCreateAccountBudgetReq{CreateAccountBudgetRequest: r}.ConvertToModel()
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
	if err := e.accountBudgetService.CreateAccountBudget(ctx, req); err != nil {
		return res, err
	}

	return res, nil
}
