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

// UpdateTag редактирует подкатегорию
func (s *TagService) UpdateTag(ctx context.Context, fields model.UpdateTagReq) error {
	ctx, span := tracer.Start(ctx, "UpdateTag")
	defer span.End()

	// Проверяем доступ пользователя к подкатегории
	if err := s.CheckAccess(ctx, fields.Necessary.UserID, []uuid.UUID{fields.ID}); err != nil {
		return err
	}

	// Получаем подкатегорию для слепка "до" в аудит-логе
	tagBefore, err := slices.FirstWithError(s.tagRepository.GetTags(ctx, model.GetTagsReq{ //nolint:exhaustruct
		IDs: []uuid.UUID{fields.ID},
	}))
	if err != nil {
		return err
	}

	return s.generalRepository.WithSyncGate(ctx, fields.Necessary.UserID, fields.Necessary.DeviceID, s.userService, s.auditLogService, func(ctxTx context.Context) (uint32, error) {

		// Изменяем данные подкатегории
		if err := s.tagRepository.UpdateTag(ctxTx, fields); err != nil {
			return 0, err
		}

		// Получаем актуальную подкатегорию из БД для слепка "после" в аудит-логе
		tagAfter, err := slices.FirstWithError(s.tagRepository.GetTags(ctxTx, model.GetTagsReq{ //nolint:exhaustruct
			IDs: []uuid.UUID{fields.ID},
		}))
		if err != nil {
			return 0, err
		}

		// Фиксируем изменение подкатегории в аудит-логе
		return s.auditLogService.TrackMutation(ctxTx, auditLogModel.TrackMutationReq{
			Entity:         auditLogEntity.Tag,
			Method:         auditLogMethod.Update,
			EntityID:       fields.ID.String(),
			Before:         tagBefore,
			After:          tagAfter,
			UserID:         fields.Necessary.UserID,
			DeviceID:       fields.Necessary.DeviceID,
			AccountGroupID: &tagAfter.AccountGroupID,
		})
	})
}
