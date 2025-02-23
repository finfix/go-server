package utils

import (
	"context"

	"pkg/errors"
	"pkg/slices"

	"server/internal/services/account/model"
	"server/internal/services/account/model/accountType"
	"server/internal/services/transaction/model/transactionType"
)

func TransactionAndAccountTypesValidation(ctx context.Context, accountFrom, accountTo model.Account, tranType transactionType.Type) error {

	var accesses string
	var isAccess bool

	// Проверяем, что типы счетов выбраны правильно для этой транзакции
	switch tranType {
	case transactionType.Income:
		isAccess = accountFrom.Type == accountType.Earnings && slices.In(accountTo.Type, accountType.Regular, accountType.Debt)
		accesses = "Earnings -> [Regular, Debt]"
	case transactionType.Transfer:
		isAccess = slices.In(accountFrom.Type, accountType.Regular, accountType.Debt) && slices.In(accountTo.Type, accountType.Regular, accountType.Debt)
		accesses = "[Regular, Debt] -> [Regular, Debt]"
	case transactionType.Consumption:
		isAccess = slices.In(accountFrom.Type, accountType.Regular, accountType.Debt) && accountTo.Type == accountType.Expense
		accesses = "[Regular, Debt] -> Expense"
	case transactionType.Balancing:
		isAccess = accountFrom.Type == accountType.Balancing && slices.In(accountTo.Type, accountType.Regular, accountType.Debt)
		accesses = "Balancing -> [Regular, Debt]"
	}

	if !isAccess {
		return errors.BadRequest.New(ctx, "Неверно выбраны типы счетов",
			errors.ParamsOption(
				"TransactionType", tranType,
				"AccountFromID", accountFrom.ID,
				"AccountToID", accountTo.ID,
				"Accesses", accesses,
			),
		)
	}

	return nil
}
