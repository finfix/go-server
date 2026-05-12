package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/settings/model"
	"server/internal/utils/necessary"

	proto "github.com/finfix/go-server-grpc/proto"
)

// UpdateCurrencies обновление курсов валют
func (e *SettingsEndpoint) UpdateCurrencies(ctx context.Context, r *proto.UpdateCurrenciesRequest) (*proto.UpdateCurrenciesResponse, error) {
	res := new(proto.UpdateCurrenciesResponse)

	// Convert proto request to internal model
	req, err := model.ProtoUpdateCurrenciesReq{UpdateCurrenciesRequest: r}.ConvertToModel()
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
	if err := e.settingsService.UpdateCurrencies(ctx, req); err != nil {
		return res, err
	}

	return res, nil
}
