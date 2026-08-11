package service

import (
	"context"
	"encoding/json"

	"server/internal/modules/auditLog/model"
	repoModel "server/internal/modules/auditLog/repository/model"
	"server/internal/utils/errors"
)

// TrackMutation фиксирует в аудит-логе изменяющее действие пользователя со слепками сущности до и после
// изменения, возвращая идентификатор созданной записи аудит-лога
func (s *AuditLogService) TrackMutation(ctx context.Context, req model.TrackMutationReq) (uint32, error) {
	ctx, span := tracer.Start(ctx, "TrackMutation")
	defer span.End()

	// Проверяем корректность названия сущности
	if err := req.Entity.Validate(ctx); err != nil {
		return 0, err
	}

	// Проверяем корректность метода изменения
	if err := req.Method.Validate(ctx); err != nil {
		return 0, err
	}

	// Сериализуем слепок сущности до изменения
	snapshotBefore, err := marshalSnapshot(ctx, req.Before)
	if err != nil {
		return 0, err
	}

	// Сериализуем слепок сущности после изменения
	snapshotAfter, err := marshalSnapshot(ctx, req.After)
	if err != nil {
		return 0, err
	}

	// Создаем запись аудит-лога
	return s.auditLogRepository.CreateAuditLog(ctx, repoModel.CreateAuditLogReq{
		Entity:         req.Entity,
		Method:         req.Method,
		EntityID:       req.EntityID,
		SnapshotBefore: snapshotBefore,
		SnapshotAfter:  snapshotAfter,
		UserID:         req.UserID,
		DeviceID:       req.DeviceID,
		AccountGroupID: req.AccountGroupID,
	})
}

// marshalSnapshot сериализует слепок сущности в JSON, возвращая nil для отсутствующего слепка
func marshalSnapshot(ctx context.Context, snapshot any) ([]byte, error) {
	if snapshot == nil {
		return nil, nil
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		return nil, errors.InternalServer.Wrap(err).WithContextParams(ctx)
	}
	return data, nil
}
