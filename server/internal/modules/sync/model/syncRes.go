package model

import (
	"github.com/google/uuid"

	"github.com/finfix/go-server-grpc/proto"

	accountModel "server/internal/modules/account/model"
	accountBudgetModel "server/internal/modules/accountBudget/model"
	accountGroupModel "server/internal/modules/accountGroup/model"
	pendingLinkedTransferModel "server/internal/modules/pendingLinkedTransfer/model"
	settingsModel "server/internal/modules/settings/model"
	tagModel "server/internal/modules/tag/model"
	transactionModel "server/internal/modules/transaction/model"
	userModel "server/internal/modules/user/model"
)

// SyncRes - результат синхронизации: изменённые и удалённые объекты по каждой сущности
type SyncRes struct {
	PendingCheckpoint uint32     // Новый чекпоинт, который клиент должен подтвердить вызовом ConfirmSync после применения изменений
	PendingSyncToken  *uuid.UUID // Опорный токен этого ответа для ConfirmSync, nil если изменений не было
	HasChanges        bool       // true, если в ответе есть изменения

	ChangedTransactions   []transactionModel.Transaction // Созданные/изменённые транзакции
	DeletedTransactionIDs []uuid.UUID                    // Удалённые транзакции

	ChangedAccounts   []accountModel.Account // Созданные/изменённые счета
	DeletedAccountIDs []uuid.UUID            // Удалённые счета

	ChangedAccountGroups   []accountGroupModel.AccountGroup // Созданные/изменённые группы счетов
	DeletedAccountGroupIDs []uuid.UUID                      // Удалённые группы счетов

	ChangedTags   []tagModel.Tag // Созданные/изменённые подкатегории
	DeletedTagIDs []uuid.UUID    // Удалённые подкатегории

	ChangedAccountBudgets []accountBudgetModel.AccountBudget // Созданные версии бюджета счетов (версии не удаляются)

	ChangedUser *userModel.User // Измененные данные текущего пользователя, nil если не менялись

	ChangedCurrencies []settingsModel.Currency // Созданные/изменённые валюты (глобальные, без привязки к группе счетов)

	ChangedPendingLinkedTransfers   []pendingLinkedTransferModel.PendingLinkedTransfer // Созданные/изменённые переносы
	DeletedPendingLinkedTransferIDs []uuid.UUID                                        // Удалённые переносы
}

// ConvertToProto преобразует SyncRes в proto-формат
func (s *SyncRes) ConvertToProto() (*proto.SyncResponse, error) {
	res := &proto.SyncResponse{ //nolint:exhaustruct
		PendingCheckpoint: s.PendingCheckpoint,
		HasChanges:        s.HasChanges,
	}
	if s.PendingSyncToken != nil {
		res.PendingSyncToken = s.PendingSyncToken[:]
	}

	for _, transaction := range s.ChangedTransactions {
		protoTransaction, err := transaction.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedTransactions = append(res.ChangedTransactions, protoTransaction)
	}
	for _, id := range s.DeletedTransactionIDs {
		res.DeletedTransactionIDs = append(res.DeletedTransactionIDs, id[:])
	}

	for _, account := range s.ChangedAccounts {
		protoAccount, err := account.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedAccounts = append(res.ChangedAccounts, protoAccount)
	}
	for _, id := range s.DeletedAccountIDs {
		res.DeletedAccountIDs = append(res.DeletedAccountIDs, id[:])
	}

	for _, accountGroup := range s.ChangedAccountGroups {
		protoAccountGroup, err := accountGroup.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedAccountGroups = append(res.ChangedAccountGroups, protoAccountGroup)
	}
	for _, id := range s.DeletedAccountGroupIDs {
		res.DeletedAccountGroupIDs = append(res.DeletedAccountGroupIDs, id[:])
	}

	for _, tag := range s.ChangedTags {
		protoTag, err := tag.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedTags = append(res.ChangedTags, protoTag)
	}
	for _, id := range s.DeletedTagIDs {
		res.DeletedTagIDs = append(res.DeletedTagIDs, id[:])
	}

	for _, budget := range s.ChangedAccountBudgets {
		res.ChangedAccountBudgets = append(res.ChangedAccountBudgets, budget.ConvertToProto())
	}

	if s.ChangedUser != nil {
		protoUser, err := s.ChangedUser.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedUser = protoUser
	}

	for _, currency := range s.ChangedCurrencies {
		protoCurrency, err := currency.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedCurrencies = append(res.ChangedCurrencies, protoCurrency)
	}

	for _, transfer := range s.ChangedPendingLinkedTransfers {
		protoTransfer, err := transfer.ConvertToProto()
		if err != nil {
			return nil, err
		}
		res.ChangedPendingLinkedTransfers = append(res.ChangedPendingLinkedTransfers, protoTransfer)
	}
	for _, id := range s.DeletedPendingLinkedTransferIDs {
		res.DeletedPendingLinkedTransferIDs = append(res.DeletedPendingLinkedTransferIDs, id[:])
	}

	return res, nil
}
