package service

import (
	"time"

	accountBudgetModel "server/internal/modules/accountBudget/model"
	"server/internal/modules/statistic/model"
)

// budgetGetReq — запрос версий бюджетов выбранных групп до конца периода
func budgetGetReq(req model.GetStatisticReq, data dataset, dateTo time.Time) accountBudgetModel.GetAccountBudgetsReq {
	end := dateTo.AddDate(0, 0, 1)
	return accountBudgetModel.GetAccountBudgetsReq{Necessary: req.Necessary, AccountGroupIDs: data.accountGroupIDs, DateTo: &end} //nolint:exhaustruct
}
