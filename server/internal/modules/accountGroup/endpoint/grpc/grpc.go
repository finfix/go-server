package grpc

import (
	"context"

	"server/internal/modules/accountGroup/model"
	accountGroupService "server/internal/modules/accountGroup/service"

	proto "github.com/finfix/go-server-grpc/proto"
)

var _ AccountGroupService = new(accountGroupService.AccountGroupService)

type AccountGroupService interface {
	CreateAccountGroup(context.Context, model.CreateAccountGroupReq) (model.CreateAccountGroupRes, error)
	GetAccountGroups(context.Context, model.GetAccountGroupsReq) ([]model.AccountGroup, error)
	UpdateAccountGroup(context.Context, model.UpdateAccountGroupReq) error
	DeleteAccountGroup(context.Context, model.DeleteAccountGroupReq) error
}

var _ proto.AccountGroupEndpointServer = new(AccountGroupEndpoint)

type AccountGroupEndpoint struct {
	proto.UnsafeAccountGroupEndpointServer
	accountGroupService AccountGroupService
}

func NewAccountGroupEndpoint(accountGroupService AccountGroupService) *AccountGroupEndpoint {
	return &AccountGroupEndpoint{
		accountGroupService: accountGroupService,
	}
}
