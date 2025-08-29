package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/accountGroup/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// CreateAccountGroup создает группу счетов
func (e *AccountGroupEndpoint) CreateAccountGroup(ctx context.Context, r *proto.CreateAccountGroupRequest) (*proto.CreateAccountGroupResponse, error) {
	res := &proto.CreateAccountGroupResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoCreateAccountGroupReq{CreateAccountGroupRequest: r}.ConvertToModel()
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
	if _, err = e.accountGroupService.CreateAccountGroup(ctx, req); err != nil {
		return res, err
	}

	// Convert response to proto and return
	return res, nil
}
