package service

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	accountModel "server/internal/modules/account/model"
	transactionModel "server/internal/modules/transaction/model"
)

// dataset — все данные выбранных групп счетов, нужные для расчета
type dataset struct {
	accountGroupIDs   []uuid.UUID                        // Выбранные группы счетов
	accounts          map[uuid.UUID]accountModel.Account // Счета по идентификатору
	transactions      []transactionModel.Transaction     // Транзакции групп, от новых к старым
	tagsByTransaction map[uuid.UUID][]uuid.UUID          // Теги каждой транзакции
	rates             map[string]decimal.Decimal         // Текущие курсы валют — запасной вариант без истории
	rateHistory       map[string][]ratePoint             // История курсов по валютам, по возрастанию момента
	approximate       map[string]bool                    // Валюты, пересчитанные приблизительно: курса до момента не было
	currency          string                             // Валюта ответа
}
