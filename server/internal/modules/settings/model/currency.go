package model

import (
	"github.com/finfix/go-server-grpc/proto"
	"github.com/shopspring/decimal"
)

type Currency struct {
	Slug   string          `json:"isoCode" db:"Slug"`  // Сигнатура валюты
	Name   string          `json:"name" db:"name"`     // Название валюты
	Symbol string          `json:"symbol" db:"symbol"` // Символ валюты
	Rate   decimal.Decimal `json:"rate" db:"rate"`     // Курс валюты
}

func (c *Currency) ConvertToProto() (res *proto.Currency, err error) {
	return &proto.Currency{
		IsoCode: c.Slug,
		Name:    c.Name,
		Symbol:  c.Symbol,
		Rate:    c.Rate.InexactFloat64(),
	}, nil
}
