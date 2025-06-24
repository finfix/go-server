package service

import (
	"context"

	"server/internal/services/tag/model"
)

func (s *TagService) GetTagsToTransactions(ctx context.Context, req model.GetTagsToTransactionsReq) (res []model.TagToTransaction, err error) {
	ctx, span := tracer.Start(ctx, "GetTagsToTransactions")
	defer span.End()

	// Получаем доступные группы счетов
	req.AccountGroupIDs, err = s.userService.GetAccessedAccountGroups(ctx, req.Necessary.UserID)
	if err != nil {
		return nil, err
	}

	// Получаем все связи между подкатегориями и транзакциями
	return s.tagRepository.GetTagsToTransactions(ctx, req)
}
