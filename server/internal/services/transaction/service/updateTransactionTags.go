package service

import (
	"context"

	"pkg/slices"

	"server/internal/services/tag/model"

	"github.com/google/uuid"
)

func (s *TransactionService) updateTransactionTags(ctx context.Context, userID, transactionID uuid.UUID, tagIDs []uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "UpdateTransactionTags")
	defer span.End()

	// Проверяем доступ пользователя к тегам
	if len(tagIDs) > 0 {
		if err := s.tagService.CheckAccess(ctx, userID, tagIDs); err != nil {
			return err
		}
	}

	// Получаем все теги, привязанные к транзакции
	transactionTags, err := s.tagRepository.GetTagsToTransactions(ctx, model.GetTagsToTransactionsReq{ //nolint:exhaustruct
		TransactionIDs: []uuid.UUID{transactionID},
	})
	if err != nil {
		return err
	}

	existTagIDs := slices.GetFields(transactionTags, func(tag model.TagToTransaction) uuid.UUID { return tag.TagID })

	toDelete, toCreate := slices.JoinExclusive(existTagIDs, tagIDs)

	if len(toDelete) > 0 {
		// Удаляем связи тегов с транзакцией
		if err = s.tagRepository.UnlinkTagsFromTransaction(ctx, toDelete, transactionID); err != nil {
			return err
		}
	}
	if len(toCreate) > 0 {
		// Создаем связи тегов с транзакцией
		if err = s.tagRepository.LinkTagsToTransaction(ctx, toCreate, transactionID); err != nil {
			return err
		}
	}
	return nil
}
