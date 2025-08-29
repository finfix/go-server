package utils

import (
	"github.com/google/uuid"

	"pkg/pointer"

	"server/internal/services/account/model"
	accountRepoModel "server/internal/services/account/repository/model"
)

func HandleVisible(
	repoUpdateReqs map[uuid.UUID]accountRepoModel.UpdateAccountReq,
	mainAccount model.Account,
	childrenAccounts []model.Account,
	parentAccount *model.Account,
) map[uuid.UUID]accountRepoModel.UpdateAccountReq {

	// Если значение родительского счета отрицательное, а у дочернего счета положительное
	if parentAccount != nil && !parentAccount.Visible && mainAccount.Visible {
		// То возникает логическая ошибка, поэтому родительский счет становится подсчитываемым
		requestParentAccount := repoUpdateReqs[parentAccount.ID]

		requestParentAccount.Visible = pointer.Pointer(true)

		repoUpdateReqs[parentAccount.ID] = requestParentAccount
	}

	// Если видимость счета отрицательная, а подсчитываемость положительного
	if !mainAccount.Visible && mainAccount.AccountingInHeader {
		// То возникает логическая ошибка, поэтому подсчитываемость деактивируется
		requestMainAccount := repoUpdateReqs[mainAccount.ID]

		requestMainAccount.AccountingInHeader = pointer.Pointer(false)

		repoUpdateReqs[mainAccount.ID] = requestMainAccount

	}

	// Если значения родительского счета меняется, то значения дочерних счетов меняются на такое же
	for _, childAccount := range childrenAccounts {
		requestChildAccount := repoUpdateReqs[childAccount.ID]

		requestChildAccount.Visible = pointer.Pointer(mainAccount.Visible)

		if !childAccount.Visible && childAccount.AccountingInHeader {
			requestChildAccount.AccountingInHeader = pointer.Pointer(false)
		}

		repoUpdateReqs[childAccount.ID] = requestChildAccount
	}

	return repoUpdateReqs
}
