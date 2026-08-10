package service

import (
	"context"

	"pkg/slices"

	"server/internal/enum/auditLogEntity"
	"server/internal/enum/auditLogMethod"
	auditLogModel "server/internal/modules/auditLog/model"
	"server/internal/modules/tag/model"

	"github.com/google/uuid"
)

// DeleteTag удаляет подкатегорию
func (s *TagService) DeleteTag(ctx context.Context, req model.DeleteTagReq) error {
	ctx, span := tracer.Start(ctx, "DeleteTag")
	defer span.End()

	// Проверяем доступ пользователя к подкатегории
	if err := s.CheckAccess(ctx, req.Necessary.UserID, []uuid.UUID{req.ID}); err != nil {
		return err
	}

	// Получаем подкатегорию для слепка "до" в аудит-логе
	tagBefore, err := slices.FirstWithError(s.tagRepository.GetTags(ctx, model.GetTagsReq{ //nolint:exhaustruct
		IDs: []uuid.UUID{req.ID},
	}))
	if err != nil {
		return err
	}

	return s.generalRepository.WithinTransaction(ctx, func(ctxTx context.Context) error {

		// Удаляем подкатегорию
		if err := s.tagRepository.DeleteTag(ctxTx, req.ID, req.Necessary.UserID); err != nil {
			return err
		}

		// Фиксируем удаление подкатегории в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:         auditLogEntity.Tag,
			Method:         auditLogMethod.Delete,
			EntityID:       req.ID.String(),
			Before:         tagBefore,
			After:          nil,
			UserID:         req.Necessary.UserID,
			DeviceID:       req.Necessary.DeviceID,
			AccountGroupID: &tagBefore.AccountGroupID,
		})
	})
}
