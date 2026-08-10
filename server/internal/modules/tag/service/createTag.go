package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/tag/model"
)

// CreateTag создает новую подкатегорию
func (s *TagService) CreateTag(ctx context.Context, tag model.CreateTagReq) (model.CreateTagRes, error) {
	ctx, span := tracer.Start(ctx, "CreateTag")
	defer span.End()

	// Проверяем доступ пользователя к группам счетов
	if err := s.accountGroupService.CheckAccess(ctx, tag.Necessary.UserID, []uuid.UUID{tag.AccountGroupID}); err != nil {
		return model.CreateTagRes{}, err
	}

	err := s.generalRepository.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Создаем подкатегорию
		if err := s.tagRepository.CreateTag(ctxTx, tag.ConvertToRepoReq()); err != nil {
			return err
		}

		// Получаем созданную подкатегорию из БД для слепка "после" в аудит-логе
		tagAfter, err := slices.FirstWithError(s.tagRepository.GetTags(ctxTx, model.GetTagsReq{ //nolint:exhaustruct
			IDs: []uuid.UUID{tag.ID},
		}))
		if err != nil {
			return err
		}

		// Фиксируем создание подкатегории в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:   auditLogEntity.Tag,
			Method:   auditLogMethod.Create,
			EntityID: tag.ID.String(),
			Before:   nil,
			After:    tagAfter,
			UserID:   tag.Necessary.UserID,
			DeviceID: tag.Necessary.DeviceID,
		})
	})
	if err != nil {
		return model.CreateTagRes{}, err
	}

	return model.CreateTagRes{}, nil
}
