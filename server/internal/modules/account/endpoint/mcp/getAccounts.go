package mcp

import (
	"context"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/shopspring/decimal"

	"server/internal/enum/accountType"
	"server/internal/modules/account/model"
	mcpTransport "server/internal/utils/mcp"
)

// GetAccountsIn — аргументы get_accounts
type GetAccountsIn struct {
	AccountGroupIDs []string `json:"accountGroupIDs,omitempty" jsonschema:"Группы счетов, пусто — все доступные"`
	Type            string   `json:"type,omitempty" jsonschema:"Тип счета: regular, debt, expense, earnings или balancing; пусто — все"`
	DateFrom        *string  `json:"dateFrom,omitempty" jsonschema:"Начало периода для remainder расходов и доходов, yyyy-MM-dd; по умолчанию первое число текущего месяца"`
	DateTo          *string  `json:"dateTo,omitempty" jsonschema:"Конец периода включительно, yyyy-MM-dd; по умолчанию сегодня"`
}

// GetAccountsOut — результат get_accounts
type GetAccountsOut struct {
	Accounts []AccountView `json:"accounts"` // Счета
}

// AccountView — счет в ответе инструмента
type AccountView struct {
	ID                 string          `json:"id"`                        // Идентификатор
	Name               string          `json:"name"`                      // Название
	Type               string          `json:"type"`                      // Тип
	Currency           string          `json:"currency"`                  // Валюта
	Remainder          decimal.Decimal `json:"remainder"`                 // Остаток или сумма за период
	Visible            bool            `json:"visible"`                   // Показывается ли в приложении
	AccountingInHeader bool            `json:"accountingInHeader"`        // Учитывается ли в шапке приложения
	AccountGroupID     string          `json:"accountGroupID"`            // Группа счетов
	ParentAccountID    *string         `json:"parentAccountID,omitempty"` // Родительский счет
	IsParent           bool            `json:"isParent"`                  // Родитель, в которого сворачиваются дочерние
	LinkedAccountID    *string         `json:"linkedAccountID,omitempty"` // Связанный счет-мост другой группы
}

// GetAccounts справочник счетов
func (t *AccountTools) GetAccounts(ctx context.Context, _ *sdk.CallToolRequest, in GetAccountsIn) (*sdk.CallToolResult, GetAccountsOut, error) {
	var out GetAccountsOut

	// Собираем запрос из аргументов и контекста
	necessary, err := mcpTransport.Necessary(ctx)
	if err != nil {
		return nil, out, err
	}
	accountGroupIDs, err := mcpTransport.ParseUUIDs("accountGroupIDs", in.AccountGroupIDs)
	if err != nil {
		return nil, out, err
	}
	dateFrom, err := mcpTransport.ParseOptionalDate("dateFrom", in.DateFrom)
	if err != nil {
		return nil, out, err
	}
	dateTo, err := mcpTransport.ParseOptionalDate("dateTo", in.DateTo)
	if err != nil {
		return nil, out, err
	}
	now := time.Now().UTC()
	if dateFrom == nil {
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		dateFrom = &monthStart
	}
	if dateTo == nil {
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		dateTo = &today
	}
	req := model.GetAccountsReq{Necessary: necessary, AccountGroupIDs: accountGroupIDs, DateFrom: dateFrom, DateTo: dateTo} //nolint:exhaustruct
	if in.Type != "" {
		accountTypeValue := accountType.AccountType(in.Type)
		if err = accountTypeValue.Validate(ctx); err != nil {
			return nil, out, mcpTransport.ToolError(err)
		}
		req.Type = &accountTypeValue
	}

	// Вызываем метод сервиса
	accounts, err := t.accountService.GetAccounts(ctx, req)
	if err != nil {
		return nil, out, mcpTransport.ToolError(err)
	}

	// Конвертируем ответ
	out.Accounts = make([]AccountView, 0, len(accounts))
	for _, account := range accounts {
		out.Accounts = append(out.Accounts, AccountView{
			ID:                 account.ID.String(),
			Name:               account.Name,
			Type:               string(account.Type),
			Currency:           account.Currency,
			Remainder:          account.Remainder,
			Visible:            account.Visible,
			AccountingInHeader: account.AccountingInHeader,
			AccountGroupID:     account.AccountGroupID.String(),
			ParentAccountID:    mcpTransport.UUIDString(account.ParentAccountID),
			IsParent:           account.IsParent,
			LinkedAccountID:    mcpTransport.UUIDString(account.LinkedAccountID),
		})
	}
	return nil, out, nil
}
