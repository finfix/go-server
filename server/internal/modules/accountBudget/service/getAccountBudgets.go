package service

import (
	"context"

	accountBudgetModel "server/internal/modules/accountBudget/model"
	"server/internal/modules/accountBudget/repository"
)

// GetAccountBudgets возвращает все версии бюджета по всем счетам всех доступных пользователю групп счетов
// (или по конкретным группам, если они переданы) - без объединения в "актуальную" версию
func (s *AccountBudgetService) GetAccountBudgets(ctx context.Context, req accountBudgetModel.GetAccountBudgetsReq) ([]accountBudgetModel.AccountBudget, error) {
	ctx, span := tracer.Start(ctx, "GetAccountBudgets")
	defer span.End()

	accountGroupIDs := req.AccountGroupIDs

	if len(accountGroupIDs) != 0 {

		// Проверяем доступ пользователя к переданным группам счетов
		if err := s.accountGroupService.CheckAccess(ctx, req.Necessary.UserID, accountGroupIDs); err != nil {
			return nil, err
		}
	} else {

		// Получаем все доступные пользователю группы счетов
		var err error
		accountGroupIDs, err = s.userToAccountGroupService.GetAccessedAccountGroups(ctx, req.Necessary.UserID)
		if err != nil {
			return nil, err
		}
	}

	return s.accountBudgetRepository.GetAccountBudgets(ctx, repository.GetAccountBudgetsReq{
		AccountGroupIDs: accountGroupIDs,
		DateFrom:        req.DateFrom,
		DateTo:          req.DateTo,
	})
}
