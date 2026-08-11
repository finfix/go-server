package service

import (
	"context"

	"github.com/google/uuid"

	accountBudgetModel "server/internal/modules/accountBudget/model"
	"server/internal/modules/accountBudget/repository"
)

// GetAccountBudgetsByIDs возвращает версии бюджета по их идентификаторам без проверки доступа - используется
// там, где видимость уже установлена по построению (например, синхронизацией по идентификаторам из аудит-лога)
func (s *AccountBudgetService) GetAccountBudgetsByIDs(ctx context.Context, ids []uuid.UUID) ([]accountBudgetModel.AccountBudget, error) {
	ctx, span := tracer.Start(ctx, "GetAccountBudgetsByIDs")
	defer span.End()

	if len(ids) == 0 {
		return nil, nil
	}

	return s.accountBudgetRepository.GetAccountBudgets(ctx, repository.GetAccountBudgetsReq{ //nolint:exhaustruct
		IDs: ids,
	})
}
