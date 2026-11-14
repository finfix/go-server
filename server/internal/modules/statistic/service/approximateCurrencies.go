package service

import "slices"

// approximateCurrencies — валюты, пересчитанные по ближайшему курсу после момента операции, по алфавиту
func (d dataset) approximateCurrencies() []string {
	result := make([]string, 0, len(d.approximate))
	for currency := range d.approximate {
		result = append(result, currency)
	}
	slices.Sort(result)
	return result
}
