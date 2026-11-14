package mcp

import (
	"github.com/shopspring/decimal"

	"server/internal/enum/transactionType"
	"server/internal/modules/statistic/model"
	mcpTransport "server/internal/utils/mcp"
)

// FiltersIn — фильтры транзакций, общие для инструментов статистики
type FiltersIn struct {
	SearchText         string           `json:"searchText,omitempty" jsonschema:"Подстрока заметки без учета регистра"`
	DateFrom           *string          `json:"dateFrom,omitempty" jsonschema:"Дата начала включительно, yyyy-MM-dd"`
	DateTo             *string          `json:"dateTo,omitempty" jsonschema:"Дата конца включительно, yyyy-MM-dd"`
	Types              []string         `json:"types,omitempty" jsonschema:"Типы транзакций: consumption, income, transfer, balancing"`
	Currencies         []string         `json:"currencies,omitempty" jsonschema:"Валюты счета списания или пополнения, например RUB"`
	AccountIDs         []string         `json:"accountIDs,omitempty" jsonschema:"Счета вместе с дочерними: транзакция списывает с них или пополняет их"`
	ExcludedAccountIDs []string         `json:"excludedAccountIDs,omitempty" jsonschema:"Исключить транзакции этих счетов вместе с дочерними"`
	TagIDs             []string         `json:"tagIDs,omitempty" jsonschema:"Транзакции с любым из тегов"`
	ExcludedTagIDs     []string         `json:"excludedTagIDs,omitempty" jsonschema:"Исключить транзакции с любым из тегов"`
	AccountGroupIDs    []string         `json:"accountGroupIDs,omitempty" jsonschema:"Группы счетов, пусто — все доступные"`
	AmountMin          *decimal.Decimal `json:"amountMin,omitempty" jsonschema:"Минимальная сумма в валюте счета: подходит, если в диапазоне сумма списания или пополнения"`
	AmountMax          *decimal.Decimal `json:"amountMax,omitempty" jsonschema:"Максимальная сумма в валюте счета"`
}

// ConvertToModel разбирает идентификаторы и даты фильтров
func (in FiltersIn) ConvertToModel() (filters model.Filters, err error) {
	filters = model.Filters{SearchText: in.SearchText, Currencies: in.Currencies, AmountMin: in.AmountMin, AmountMax: in.AmountMax} //nolint:exhaustruct
	if filters.DateFrom, err = mcpTransport.ParseOptionalDate("dateFrom", in.DateFrom); err != nil {
		return filters, err
	}
	if filters.DateTo, err = mcpTransport.ParseOptionalDate("dateTo", in.DateTo); err != nil {
		return filters, err
	}
	for _, item := range in.Types {
		filters.Types = append(filters.Types, transactionType.TransactionType(item))
	}
	if filters.AccountIDs, err = mcpTransport.ParseUUIDs("accountIDs", in.AccountIDs); err != nil {
		return filters, err
	}
	if filters.ExcludedAccountIDs, err = mcpTransport.ParseUUIDs("excludedAccountIDs", in.ExcludedAccountIDs); err != nil {
		return filters, err
	}
	if filters.TagIDs, err = mcpTransport.ParseUUIDs("tagIDs", in.TagIDs); err != nil {
		return filters, err
	}
	if filters.ExcludedTagIDs, err = mcpTransport.ParseUUIDs("excludedTagIDs", in.ExcludedTagIDs); err != nil {
		return filters, err
	}
	if filters.AccountGroupIDs, err = mcpTransport.ParseUUIDs("accountGroupIDs", in.AccountGroupIDs); err != nil {
		return filters, err
	}
	return filters, nil
}
