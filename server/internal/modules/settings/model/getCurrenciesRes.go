package model

import "github.com/finfix/go-server-grpc/proto"

type GetCurrenciesRes struct {
	Currencies []Currency
}

func (s *GetCurrenciesRes) ConvertToProto() (res *proto.GetCurrenciesResponse, err error) {

	protoCurrencies := make([]*proto.Currency, 0, len(s.Currencies))
	for _, currency := range s.Currencies {
		protoCurrency, err := currency.ConvertToProto()
		if err != nil {
			return res, err
		}
		protoCurrencies = append(protoCurrencies, protoCurrency)
	}

	return &proto.GetCurrenciesResponse{
		Error:      nil,
		Currencies: protoCurrencies,
	}, nil
}
