package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/accountGroup/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// DeleteAccountGroup удаляет группу счетов
func (e *AccountGroupEndpoint) DeleteAccountGroup(ctx context.Context, r *proto.DeleteAccountGroupRequest) (*proto.DeleteAccountGroupResponse, error) {
	res := &proto.DeleteAccountGroupResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoDeleteAccountGroupReq{DeleteAccountGroupRequest: r}.ConvertToModel()
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
	err = e.accountGroupService.DeleteAccountGroup(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}
