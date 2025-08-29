package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/accountGroup/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// UpdateAccountGroup обновляет группу счетов
func (e *AccountGroupEndpoint) UpdateAccountGroup(ctx context.Context, r *proto.UpdateAccountGroupRequest) (*proto.UpdateAccountGroupResponse, error) {
	res := &proto.UpdateAccountGroupResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoUpdateAccountGroupReq{UpdateAccountGroupRequest: r}.ConvertToModel()
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
	err = e.accountGroupService.UpdateAccountGroup(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}
