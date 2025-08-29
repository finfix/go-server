package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/account/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// GetAccounts получает счета по фильтрам
func (e *AccountEndpoint) GetAccounts(ctx context.Context, r *proto.GetAccountsRequest) (*proto.GetAccountsResponse, error) {
	res := &proto.GetAccountsResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoGetAccountsReq{GetAccountsRequest: r}.ConvertToModel()
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
	accounts, err := e.accountService.GetAccounts(ctx, req)
	if err != nil {
		return res, err
	}

	// Convert accounts to proto format
	protoAccounts := make([]*proto.Account, 0, len(accounts))
	for _, account := range accounts {
		protoAccount, err := account.ConvertToProto()
		if err != nil {
			return res, err
		}
		protoAccounts = append(protoAccounts, protoAccount)
	}

	res.Accounts = protoAccounts
	return res, nil
}
