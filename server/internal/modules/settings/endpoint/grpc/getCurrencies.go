package grpc

import (
	"context"

	"github.com/finfix/go-server-grpc/proto"
)

// GetCurrencies получение списка валют
func (e *SettingsEndpoint) GetCurrencies(ctx context.Context, r *proto.GetCurrenciesRequest) (*proto.GetCurrenciesResponse, error) {
	res := new(proto.GetCurrenciesResponse)

	// Call service method
	_res, err := e.settingsService.GetCurrencies(ctx)
	if err != nil {
		return res, err
	}

	return _res.ConvertToProto()
}
