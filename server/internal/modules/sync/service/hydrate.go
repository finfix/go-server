package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	accountModel "server/internal/modules/account/model"
	accountGroupModel "server/internal/modules/accountGroup/model"
	auditLogModel "server/internal/modules/auditLog/model"
	pendingLinkedTransferModel "server/internal/modules/pendingLinkedTransfer/model"
	"server/internal/modules/sync/model"
	tagModel "server/internal/modules/tag/model"
	transactionModel "server/internal/modules/transaction/model"
	userModel "server/internal/modules/user/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

// hydrate раскладывает последние события аудит-лога по сущностям: удалённые объекты собираются
// в списки идентификаторов напрямую, а созданные/изменённые - довычитываются актуальными из своих
// доменных сервисов (а не восстанавливаются из слепка аудит-лога), чтобы не отдавать клиенту
// устаревшие вычисляемые поля (например, остаток счета)
func (s *SyncService) hydrate(ctx context.Context, userID uuid.UUID, latest map[entityKey]auditLogModel.AuditLog, res *model.SyncRes) error {
	var transactionIDs, accountIDs, accountGroupIDs, tagIDs, accountBudgetIDs, pendingLinkedTransferIDs []uuid.UUID
	var hasUserChange, hasCurrencyChange bool

	for key, auditLog := range latest {
		entity := auditLogEntity.AuditLogEntity(key.entity)

		// Currency - единственная сущность, чей entity_id не является UUID (это slug валюты,
		// например "usd"), поэтому обрабатываем её до попытки парсинга UUID
		if entity == auditLogEntity.Currency {
			// Валюты глобальны и не удаляются - метод всегда create/update
			hasCurrencyChange = true
			continue
		}

		id, err := uuid.Parse(key.entityID)
		if err != nil {
			return errors.InternalServer.Wrap(err).WithContextParams(ctx)
		}

		switch entity {
		case auditLogEntity.Transaction:
			if auditLog.Method == auditLogMethod.Delete {
				res.DeletedTransactionIDs = append(res.DeletedTransactionIDs, id)
			} else {
				transactionIDs = append(transactionIDs, id)
			}

		case auditLogEntity.Account:
			if auditLog.Method == auditLogMethod.Delete {
				res.DeletedAccountIDs = append(res.DeletedAccountIDs, id)
			} else {
				accountIDs = append(accountIDs, id)
			}

		case auditLogEntity.AccountGroup:
			if auditLog.Method == auditLogMethod.Delete {
				res.DeletedAccountGroupIDs = append(res.DeletedAccountGroupIDs, id)
			} else {
				accountGroupIDs = append(accountGroupIDs, id)
			}

		case auditLogEntity.Tag:
			if auditLog.Method == auditLogMethod.Delete {
				res.DeletedTagIDs = append(res.DeletedTagIDs, id)
			} else {
				tagIDs = append(tagIDs, id)
			}

		case auditLogEntity.AccountBudget:
			// Версии бюджета счетов неизменяемы и не удаляются - метод всегда create
			accountBudgetIDs = append(accountBudgetIDs, id)

		case auditLogEntity.User:
			// Пользователь не удаляется - метод всегда create/update. Видна только запись самого
			// пользователя (отфильтровано на уровне GetAuditLogsSince), поэтому id всегда его собственный
			hasUserChange = true

		case auditLogEntity.PendingLinkedTransfer:
			// Переносы не удаляются - метод всегда create/update (см. status)
			pendingLinkedTransferIDs = append(pendingLinkedTransferIDs, id)
		}
	}

	necessaryInfo := necessary.NecessaryUserInformation{UserID: userID}
	now := time.Now()

	var err error

	if len(transactionIDs) != 0 {
		res.ChangedTransactions, err = s.transactionService.GetTransactions(ctx, transactionModel.GetTransactionsReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       transactionIDs,
		})
		if err != nil {
			return err
		}
	}

	if len(accountIDs) != 0 {
		res.ChangedAccounts, err = s.accountService.GetAccounts(ctx, accountModel.GetAccountsReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       accountIDs,
			DateFrom:  &now,
			DateTo:    &now,
		})
		if err != nil {
			return err
		}
	}

	if len(accountGroupIDs) != 0 {
		res.ChangedAccountGroups, err = s.accountGroupService.GetAccountGroups(ctx, accountGroupModel.GetAccountGroupsReq{ //nolint:exhaustruct
			AccountGroupIDs: accountGroupIDs,
		})
		if err != nil {
			return err
		}
	}

	if len(tagIDs) != 0 {
		res.ChangedTags, err = s.tagService.GetTags(ctx, tagModel.GetTagsReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       tagIDs,
		})
		if err != nil {
			return err
		}
	}

	if len(accountBudgetIDs) != 0 {
		res.ChangedAccountBudgets, err = s.accountBudgetService.GetAccountBudgetsByIDs(ctx, accountBudgetIDs)
		if err != nil {
			return err
		}
	}

	if hasUserChange {
		users, err := s.userService.GetUsers(ctx, userModel.GetUsersReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       []uuid.UUID{userID},
		})
		if err != nil {
			return err
		}
		if len(users) != 0 {
			res.ChangedUser = &users[0]
		}
	}

	if hasCurrencyChange {
		currencies, err := s.settingsService.GetCurrencies(ctx)
		if err != nil {
			return err
		}
		res.ChangedCurrencies = currencies.Currencies
	}

	if len(pendingLinkedTransferIDs) != 0 {
		res.ChangedPendingLinkedTransfers, err = s.pendingLinkedTransferService.GetPendingLinkedTransfers(ctx, pendingLinkedTransferModel.GetPendingLinkedTransfersReq{ //nolint:exhaustruct
			Necessary: necessaryInfo,
			IDs:       pendingLinkedTransferIDs,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
