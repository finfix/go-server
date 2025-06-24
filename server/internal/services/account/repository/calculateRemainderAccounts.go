package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/shopspring/decimal"

	"pkg/ddlHelper"
	"server/internal/services/account/repository/accountDDL"
	accountRepoModel "server/internal/services/account/repository/model"
	"server/internal/services/transaction/repository/transactionDDL"
)

type amountsArray struct {
	ID        uint32          `db:"id"`
	Remainder decimal.Decimal `db:"remainder"`
}

type SumTransactionsType int

const (
	SumTransactionsToAccount SumTransactionsType = iota + 1
	SumTransactionsFromAccount
)

func buildRequestForGettingSumTransactions(req accountRepoModel.CalculateRemaindersAccountsReq, mode SumTransactionsType) sq.SelectBuilder {

	var selectField, sumField string

	// Выбираем поля в зависимости от типа транзакции
	switch mode {
	case SumTransactionsToAccount:
		selectField = transactionDDL.ColumnAccountToID
		sumField = transactionDDL.ColumnAmountTo
	case SumTransactionsFromAccount:
		selectField = transactionDDL.ColumnAccountFromID
		sumField = transactionDDL.ColumnAmountFrom
	}

	// Формируем запрос
	q := sq.
		Select(
			ddlHelper.As(
				transactionDDL.WithPrefix(selectField),
				"id",
			),
			ddlHelper.As(
				ddlHelper.Coalesce(
					ddlHelper.Sum(transactionDDL.WithPrefix(sumField)),
					"0",
				),
				"remainder",
			),
		).
		From(transactionDDL.TableWithAlias).
		Join(ddlHelper.BuildJoin(
			accountDDL.TableWithAlias,
			accountDDL.WithPrefix(accountDDL.ColumnID),
			transactionDDL.WithPrefix(selectField),
		)).
		GroupBy(transactionDDL.WithPrefix(selectField))

	// Дополняем запрос фильтрами в зависимости от параметров
	if len(req.IDs) != 0 {
		q = q.Where(sq.Eq{accountDDL.WithPrefix(accountDDL.ColumnID): req.IDs})
	}
	if len(req.Types) != 0 {
		q = q.Where(sq.Eq{accountDDL.WithPrefix(accountDDL.ColumnType): req.Types})
	}
	if len(req.AccountGroupIDs) != 0 {
		q = q.Where(sq.Eq{accountDDL.WithPrefix(accountDDL.ColumnAccountGroupID): req.AccountGroupIDs})
	}
	if req.DateFrom != nil {
		q = q.Where(sq.GtOrEq{transactionDDL.WithPrefix(transactionDDL.ColumnDate): req.DateFrom})
	}
	if req.DateTo != nil {
		q = q.Where(sq.LtOrEq{transactionDDL.WithPrefix(transactionDDL.ColumnDate): req.DateTo})
	}

	return q
}

// GetSumAllTransactionsToAccount возвращает суммы всех транзакций, которые исходили из счетов
func (r *AccountRepository) GetSumAllTransactionsToAccount(ctx context.Context, req accountRepoModel.CalculateRemaindersAccountsReq) (map[uint32]decimal.Decimal, error) {
	ctx, span := tracer.Start(ctx, "GetSumAllTransactionsToAccount")
	defer span.End()

	var amountsArray []amountsArray

	// Вычисляем сумму всех транзакций из счетов, id - сумма из
	if err := r.db.Select(ctx, &amountsArray, buildRequestForGettingSumTransactions(req, SumTransactionsToAccount)); err != nil {
		return nil, err
	}

	// Формируем мапу с суммой исходящих транзакций в виде map[accountID]amountTransactions
	amountFromAccount := make(map[uint32]decimal.Decimal, len(amountsArray))
	for _, remainder := range amountsArray {
		amountFromAccount[remainder.ID] = remainder.Remainder
	}

	amountsArray = nil

	// Вычисляем сумму всех транзакций в счета, id - сумма в
	if err := r.db.Select(ctx, &amountsArray, buildRequestForGettingSumTransactions(req, SumTransactionsFromAccount)); err != nil {
		return nil, err
	}

	// Формируем мапу с суммой входящих транзакций в виде map[accountID]amountTransactions
	amountToAccount := make(map[uint32]decimal.Decimal, len(amountFromAccount))
	for _, remainder := range amountsArray {
		amountToAccount[remainder.ID] = remainder.Remainder
	}

	// Проходим по всем счетам и вычисляем остаток разницей суммы в и из счета, формируем новую мапу
	amountMapToAccountID := make(map[uint32]decimal.Decimal)
	for id := range amountToAccount {
		amountMapToAccountID[id] = amountFromAccount[id].Sub(amountToAccount[id])
	}

	// Если счета нет в списке транзакций, то добавляем его со значением -сумма из счета
	for id := range amountFromAccount {
		if _, ok := amountMapToAccountID[id]; !ok {
			amountMapToAccountID[id] = amountFromAccount[id]
		}
	}

	return amountMapToAccountID, nil
}
