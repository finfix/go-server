package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/account/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// CreateAccount создает счет в системе
func (e *AccountEndpoint) CreateAccount(ctx context.Context, r *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	res := &proto.CreateAccountResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoCreateAccountReq{CreateAccountRequest: r}.ConvertToModel()
	if err != nil {
		// TODO: Convert error to proto error format
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
	createRes, err := e.accountService.CreateAccount(ctx, req)
	if err != nil {
		return res, err
	}

	// Convert response to proto and return
	return createRes.ConvertToProto(), nil
}
