package grpc

import (
	"context"

	"server/internal/modules/accountBudget/model"
	accountBudgetService "server/internal/modules/accountBudget/service"

	"github.com/finfix/go-server-grpc/proto"
)

var _ AccountBudgetService = new(accountBudgetService.AccountBudgetService)

type AccountBudgetService interface {
	CreateAccountBudget(context.Context, model.CreateAccountBudgetReq) error
	GetAccountBudgets(context.Context, model.GetAccountBudgetsReq) ([]model.AccountBudget, error)
}

var _ proto.AccountBudgetEndpointServer = new(AccountBudgetEndpoint)

type AccountBudgetEndpoint struct {
	proto.UnsafeAccountBudgetEndpointServer
	accountBudgetService AccountBudgetService
}

func NewAccountBudgetEndpoint(accountBudgetService AccountBudgetService) *AccountBudgetEndpoint {
	return &AccountBudgetEndpoint{
		accountBudgetService: accountBudgetService,
	}
}
