package model

import (
	"server/internal/utils/necessary"
)

// GetTransactionsReq — запрос ленты транзакций
type GetTransactionsReq struct {
	Necessary necessary.NecessaryUserInformation
	Currency  string  // Валюта итогов, пусто — по умолчанию
	Filters   Filters // Фильтры
	Offset    int     // Сколько транзакций пропустить (курсор)
	Limit     int     // Сколько транзакций вернуть
}
