package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/account/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// UpdateAccount обновляет счет
func (e *AccountEndpoint) UpdateAccount(ctx context.Context, r *proto.UpdateAccountRequest) (*proto.UpdateAccountResponse, error) {
	res := &proto.UpdateAccountResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoUpdateAccountReq{UpdateAccountRequest: r}.ConvertToModel()
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
	_, err = e.accountService.UpdateAccount(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}
