package model

import (
	"github.com/finfix/go-server-grpc/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertToProto конвертирует страницу истории курсов в proto-ответ; курс — строкой полной точности
func (r GetCurrencyRatesRes) ConvertToProto() *proto.GetCurrencyRatesResponse {
	res := &proto.GetCurrencyRatesResponse{Rates: make([]*proto.CurrencyRate, 0, len(r.Rates)), Truncated: r.HasMore} //nolint:exhaustruct
	for _, rate := range r.Rates {
		res.Rates = append(res.Rates, &proto.CurrencyRate{IsoCode: rate.Currency, ValidAt: timestamppb.New(rate.ValidAt), Rate: rate.Rate.String()}) //nolint:exhaustruct
	}
	return res
}
