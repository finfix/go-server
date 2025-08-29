package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/accountGroup/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// GetAccountGroups получает список групп счетов
func (e *AccountGroupEndpoint) GetAccountGroups(ctx context.Context, r *proto.GetAccountGroupsRequest) (*proto.GetAccountGroupsResponse, error) {
	res := &proto.GetAccountGroupsResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoGetAccountGroupsReq{GetAccountGroupsRequest: r}.ConvertToModel()
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
	accountGroups, err := e.accountGroupService.GetAccountGroups(ctx, req)
	if err != nil {
		return res, err
	}

	// Convert account groups to proto format
	protoAccountGroups := make([]*proto.AccountGroup, 0, len(accountGroups))
	for _, accountGroup := range accountGroups {
		protoAccountGroup, err := accountGroup.ConvertToProto()
		if err != nil {
			return res, err
		}
		protoAccountGroups = append(protoAccountGroups, protoAccountGroup)
	}

	res.AccountGroups = protoAccountGroups
	return res, nil
}
