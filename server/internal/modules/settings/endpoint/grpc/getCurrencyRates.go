package grpc

import (
	"context"

	"github.com/finfix/go-server-grpc/proto"

	"server/internal/modules/settings/model"
)

// GetCurrencyRates история курсов валют страницами
func (e *SettingsEndpoint) GetCurrencyRates(ctx context.Context, r *proto.GetCurrencyRatesRequest) (*proto.GetCurrencyRatesResponse, error) {
	res := new(proto.GetCurrencyRatesResponse)

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoGetCurrencyRatesReq{GetCurrencyRatesRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Получаем страницу истории
	page, err := e.settingsService.GetCurrencyRatesPage(ctx, req)
	if err != nil {
		return res, err
	}
	return page.ConvertToProto(), nil
}
