package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/account/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// DeleteAccount удаляет счет
func (e *AccountEndpoint) DeleteAccount(ctx context.Context, r *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error) {
	res := &proto.DeleteAccountResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoDeleteAccountReq{DeleteAccountRequest: r}.ConvertToModel()
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
	err = e.accountService.DeleteAccount(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}
