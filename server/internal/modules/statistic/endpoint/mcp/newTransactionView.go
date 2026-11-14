package mcp

import (
	"server/internal/modules/statistic/model"
)

// NewTransactionView конвертирует транзакцию ленты в ответ инструмента
func NewTransactionView(transaction model.TransactionItem) TransactionView {
	tagIDs := make([]string, 0, len(transaction.TagIDs))
	for _, id := range transaction.TagIDs {
		tagIDs = append(tagIDs, id.String())
	}
	return TransactionView{
		ID:              transaction.ID.String(),
		Type:            string(transaction.Type),
		AccountFromID:   transaction.AccountFromID.String(),
		AccountFromName: transaction.AccountFromName,
		AmountFrom:      transaction.AmountFrom,
		CurrencyFrom:    transaction.CurrencyFrom,
		AccountToID:     transaction.AccountToID.String(),
		AccountToName:   transaction.AccountToName,
		AmountTo:        transaction.AmountTo,
		CurrencyTo:      transaction.CurrencyTo,
		Note:            transaction.Note,
		TagIDs:          tagIDs,
		AccountGroupID:  transaction.AccountGroupID.String(),
		CreatedByUserID: transaction.CreatedByUserID.String(),
	}
}
