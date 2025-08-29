package service

import (
	"context"

	"server/internal/services/tag/model"

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

	// Удаляем подкатегорию
	return s.tagRepository.DeleteTag(ctx, req.ID, req.Necessary.UserID)
}
