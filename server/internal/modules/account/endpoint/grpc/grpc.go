package grpc

import (
	"context"

	"server/internal/modules/account/model"
	accountService "server/internal/modules/account/service"

	proto "github.com/finfix/go-server-grpc/proto"
)

var _ AccountService = new(accountService.AccountService)

type AccountService interface {
	CreateAccount(context.Context, model.CreateAccountReq) (model.CreateAccountRes, error)
	GetAccounts(context.Context, model.GetAccountsReq) ([]model.Account, error)
	UpdateAccount(context.Context, model.UpdateAccountReq) (model.UpdateAccountRes, error)
	DeleteAccount(context.Context, model.DeleteAccountReq) error
}

var _ proto.AccountEndpointServer = new(AccountEndpoint)

type AccountEndpoint struct {
	proto.UnsafeAccountEndpointServer
	accountService AccountService
}

func NewAccountEndpoint(accountService AccountService) *AccountEndpoint {
	return &AccountEndpoint{
		accountService: accountService,
	}
}
