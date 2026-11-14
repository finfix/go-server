package service

import (
	"context"

	"github.com/shopspring/decimal"

	"server/internal/enum/transactionType"
	"server/internal/modules/statistic/model"
	"server/internal/utils/errors"
)

// Пределы страницы ленты
const (
	defaultTransactionsLimit = 100 // Размер страницы по умолчанию
	maxTransactionsLimit     = 500 // Наибольший размер страницы
)

// GetTransactions возвращает страницу ленты, как вкладка транзакций приложения: от новых к старым,
// по дням с суммой расходов за день, плюс итоги по всей выборке
func (s *StatisticService) GetTransactions(ctx context.Context, req model.GetTransactionsReq) (res model.TransactionsPage, err error) {
	ctx, span := tracer.Start(ctx, "GetTransactions")
	defer span.End()

	// Проверяем пагинацию
	if req.Limit == 0 {
		req.Limit = defaultTransactionsLimit
	}
	if req.Limit < 0 || req.Limit > maxTransactionsLimit || req.Offset < 0 {
		return res, errors.BadRequest.New("limit must be 1..500 and offset must be >= 0").
			WithCustomHumanText("limit — от 1 до 500, offset — неотрицательный")
	}

	// Данные выбранных групп
	data, err := s.loadDataset(ctx, req.Necessary, req.Filters.AccountGroupIDs, req.Currency)
	if err != nil {
		return res, err
	}
	res.Currency = data.currency
	res.TotalConsumption, res.TotalIncome = decimal.Zero, decimal.Zero

	// Выборка по фильтрам с итогами и расходами по дням
	filter := data.newTransactionFilter(req.Filters, false)
	dayConsumption := make(map[string]decimal.Decimal)
	var selected []int
	for i, transaction := range data.transactions {
		if !filter.inDates(transaction) || !filter.match(data, transaction) {
			continue
		}
		selected = append(selected, i)

		switch transaction.Type {
		case transactionType.Consumption:
			amount, err := data.convert(transaction.AmountTo, data.accounts[transaction.AccountToID].Currency, endOfDay(transaction.DateTransaction))
			if err != nil {
				return res, err
			}
			res.TotalConsumption = res.TotalConsumption.Add(amount)
			day := transaction.DateTransaction.Format("2006-01-02")
			dayConsumption[day] = dayConsumption[day].Add(amount)
		case transactionType.Income:
			amount, err := data.convert(transaction.AmountFrom, data.accounts[transaction.AccountFromID].Currency, endOfDay(transaction.DateTransaction))
			if err != nil {
				return res, err
			}
			res.TotalIncome = res.TotalIncome.Add(amount)
		}
	}
	res.TotalCount = len(selected)
	res.ApproximateCurrencies = data.approximateCurrencies()

	// Страница
	if req.Offset >= len(selected) {
		return res, nil
	}
	end := min(req.Offset+req.Limit, len(selected))
	if end < len(selected) {
		res.NextOffset = &end
	}

	// Раскладываем страницу по дням: транзакции уже отсортированы от новых к старым
	for _, i := range selected[req.Offset:end] {
		transaction := data.transactions[i]
		day := transaction.DateTransaction.Format("2006-01-02")
		if len(res.Days) == 0 || res.Days[len(res.Days)-1].Date.Format("2006-01-02") != day {
			res.Days = append(res.Days, model.TransactionsDay{
				Date:             startOf(transaction.DateTransaction, model.PeriodDay),
				ConsumptionTotal: dayConsumption[day].Add(decimal.Zero),
				Transactions:     nil,
			})
		}
		current := &res.Days[len(res.Days)-1]
		current.Transactions = append(current.Transactions, model.TransactionItem{
			ID:              transaction.ID,
			Type:            transaction.Type,
			Date:            transaction.DateTransaction,
			AccountFromID:   transaction.AccountFromID,
			AccountFromName: data.accounts[transaction.AccountFromID].Name,
			AmountFrom:      transaction.AmountFrom,
			CurrencyFrom:    data.accounts[transaction.AccountFromID].Currency,
			AccountToID:     transaction.AccountToID,
			AccountToName:   data.accounts[transaction.AccountToID].Name,
			AmountTo:        transaction.AmountTo,
			CurrencyTo:      data.accounts[transaction.AccountToID].Currency,
			Note:            transaction.Note,
			TagIDs:          data.tagsByTransaction[transaction.ID],
			AccountGroupID:  transaction.AccountGroupID,
			CreatedByUserID: transaction.CreatedByUserID,
		})
	}
	return res, nil
}
