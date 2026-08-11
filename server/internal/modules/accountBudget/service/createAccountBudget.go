package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	accountModel "server/internal/modules/account/model"
	accountBudgetModel "server/internal/modules/accountBudget/model"
	auditLogModel "server/internal/modules/auditLog/model"
)

// CreateAccountBudget создает новую версию бюджета счета
func (s *AccountBudgetService) CreateAccountBudget(ctx context.Context, req accountBudgetModel.CreateAccountBudgetReq) error {
	ctx, span := tracer.Start(ctx, "CreateAccountBudget")
	defer span.End()

	// Проверяем доступ пользователя к счету
	if err := s.accountService.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{req.AccountID}); err != nil {
		return err
	}

	// Получаем счет, чтобы узнать группу счетов. DateFrom/DateTo здесь не имеют смысла (остаток счета
	// не используется), но обязательны для расчета остатков счетов любого типа внутри GetAccounts
	account, err := slices.FirstWithError(s.accountService.GetAccounts(ctx, accountModel.GetAccountsReq{ //nolint:exhaustruct
		Necessary: req.Necessary,
		IDs:       []uuid.UUID{req.AccountID},
		DateFrom:  &req.EffectiveFrom,
		DateTo:    &req.EffectiveFrom,
	}))
	if err != nil {
		return err
	}

	return s.transactor.WithinTransaction(ctx, func(ctxTx context.Context) error {

		budget := accountBudgetModel.AccountBudget{
			ID:              req.ID,
			AccountID:       req.AccountID,
			AccountGroupID:  account.AccountGroupID,
			Amount:          req.Amount,
			FixedSum:        req.FixedSum,
			DaysOffset:      req.DaysOffset,
			GradualFilling:  req.GradualFilling,
			EffectiveFrom:   req.EffectiveFrom,
			CreatedByUserID: req.Necessary.UserID,
		}

		// Создаем версию бюджета
		if err := s.accountBudgetRepository.CreateAccountBudget(ctxTx, budget); err != nil {
			return err
		}

		// Фиксируем создание версии бюджета в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:         auditLogEntity.AccountBudget,
			Method:         auditLogMethod.Create,
			EntityID:       req.ID.String(),
			Before:         nil,
			After:          budget,
			UserID:         req.Necessary.UserID,
			DeviceID:       req.Necessary.DeviceID,
			AccountGroupID: &account.AccountGroupID,
		})
	})
}
