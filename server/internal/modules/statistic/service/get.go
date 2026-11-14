package service

import (
	"github.com/shopspring/decimal"

	"server/internal/modules/statistic/model"
)

// get возвращает ряд по ключу, создавая пустой
func (c seriesCalculator) get(key, name string) *model.Series {
	if series, ok := c.series[key]; ok {
		return series
	}
	series := &model.Series{Key: key, Name: name, Values: make([]decimal.Decimal, len(c.buckets))} //nolint:exhaustruct
	for i := range series.Values {
		series.Values[i] = decimal.Zero
	}
	c.series[key] = series
	return series
}
