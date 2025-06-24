package service

import (
	"context"

	"pkg/slices"
	"server/internal/utils/errors"
)

func (s *AccountGroupService) CheckAccess(ctx context.Context, userID uint32, accountGroupIDs []uint32) error {
	ctx, span := tracer.Start(ctx, "CheckAccess")
	defer span.End()

	// Получаем группы счетов, к которым есть доступ у пользователя
	accessedAccountGroupIDs, err := s.userService.GetAccessedAccountGroups(ctx, userID)
	if err != nil {
		return err
	}

	// Если доступных групп счетов нет, возвращаем ошибку
	if len(accessedAccountGroupIDs) == 0 {
		return errors.NotFound.New("Нет доступных групп счетов").
			WithContextParams(ctx).
			WithParams("UserID", userID)
	}

	// Преобразуем доступные группы счетов в map
	accessedAccountGroupIDsMap := slices.ToMap(accessedAccountGroupIDs, func(userID uint32) uint32 { return userID })

	// Проходимся по каждой доступной группе счетов
	for _, accountGroupID := range accountGroupIDs {

		// Если нет доступа к группе счетов
		if _, ok := accessedAccountGroupIDsMap[accountGroupID]; !ok {

			// Возвращаем ошибку
			return errors.Forbidden.New("Access denied").
				WithContextParams(ctx).
				WithParams(
					"UserID", userID,
					"AccountGroupID", accountGroupIDs,
				).
				WithCustomHumanText("Вы не имеете доступа к группе счетов с ID = %v", accountGroupID).
				SkipPreviousCaller()
		}
	}

	return nil

}
