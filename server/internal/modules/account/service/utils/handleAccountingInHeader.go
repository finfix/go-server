package utils

import (
	"github.com/google/uuid"

	"pkg/pointer"

	"server/internal/modules/account/model"
	accountRepoModel "server/internal/modules/account/repository/model"
)

func HandleAccountingInHeader(
	repoUpdateReqs map[uuid.UUID]accountRepoModel.UpdateAccountReq,
	mainAccount model.Account,
	childrenAccounts []model.Account,
	parentAccount *model.Account,
) map[uuid.UUID]accountRepoModel.UpdateAccountReq {

	// Если значение родительского счета отрицательное, а у дочернего счета положительное
	if parentAccount != nil && !parentAccount.AccountingInHeader && mainAccount.AccountingInHeader {
		// То возникает логическая ошибка, поэтому родительский счет становится подсчитываемым
		requestParentAccount := repoUpdateReqs[parentAccount.ID]

		requestParentAccount.AccountingInHeader = pointer.Pointer(true)

		repoUpdateReqs[parentAccount.ID] = requestParentAccount
	}

	for _, childAccount := range childrenAccounts {

		if childAccount.AccountingInHeader && !mainAccount.AccountingInHeader {
			// То возникает логическая ошибка, поэтому значение дочернего счета становится отрицательным
			requestChildAccount := repoUpdateReqs[childAccount.ID]

			requestChildAccount.AccountingInHeader = pointer.Pointer(false)

			repoUpdateReqs[childAccount.ID] = requestChildAccount
		}

	}

	return repoUpdateReqs
}
