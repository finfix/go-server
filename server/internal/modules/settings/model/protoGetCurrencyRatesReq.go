package model

import (
	"github.com/finfix/go-server-grpc/proto"

	"server/internal/utils/errors"
)

// ProtoGetCurrencyRatesReq — обертка над proto-запросом истории курсов
type ProtoGetCurrencyRatesReq struct {
	*proto.GetCurrencyRatesRequest
}

// ConvertToModel конвертирует proto-запрос во внутреннюю модель
func (p ProtoGetCurrencyRatesReq) ConvertToModel() (GetCurrencyRatesReq, error) {
	if p.GetCurrencyRatesRequest == nil {
		return GetCurrencyRatesReq{}, errors.BadRequest.New("GetCurrencyRatesRequest is required")
	}
	req := GetCurrencyRatesReq{Currencies: p.IsoCodes, Limit: int(p.Limit)} //nolint:exhaustruct
	if p.Since != nil {
		since := p.Since.AsTime()
		req.Since = &since
	}
	return req, nil
}
