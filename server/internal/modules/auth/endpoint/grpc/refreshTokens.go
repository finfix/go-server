package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/auth/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// RefreshTokens обновление токенов
func (e *AuthEndpoint) RefreshTokens(ctx context.Context, r *proto.RefreshTokensRequest) (*proto.RefreshTokensResponse, error) {
	res := new(proto.RefreshTokensResponse)

	// Convert proto request to internal model
	req, err := model.ProtoRefreshTokensReq{RefreshTokensRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return res, err
	}

	// Call service method
	refreshRes, err := e.authService.RefreshTokens(ctx, req)
	if err != nil {
		return res, err
	}

	// Convert response to proto and return
	return refreshRes.ConvertToProto()
}
