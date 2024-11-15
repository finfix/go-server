package service

import (
	"context"

	"server/internal/services/tag/model"
)

// UpdateTag редактирует подкатегорию
func (s *TagService) UpdateTag(ctx context.Context, fields model.UpdateTagReq) error {
	ctx, span := tracer.Start(ctx, "UpdateTag")
	defer span.End()

	// Проверяем доступ пользователя к подкатегории
	if err := s.CheckAccess(ctx, fields.Necessary.UserID, []uint32{fields.ID}); err != nil {
		return err
	}

	// Изменяем данные подкатегории
	return s.tagRepository.UpdateTag(ctx, fields)
}
