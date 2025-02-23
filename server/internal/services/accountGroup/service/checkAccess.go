package service

import (
	"context"

	"pkg/errors"
	"pkg/slices"
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
		return errors.NotFound.New(ctx, "Нет доступных групп счетов",
			errors.ParamsOption("UserID", userID),
		)
	}

	// Преобразуем доступные группы счетов в map
	accessedAccountGroupIDsMap := slices.ToMap(accessedAccountGroupIDs, func(userID uint32) uint32 { return userID })

	// Проходимся по каждой доступной группе счетов
	for _, accountGroupID := range accountGroupIDs {

		// Если нет доступа к группе счетов
		if _, ok := accessedAccountGroupIDsMap[accountGroupID]; !ok {

			// Возвращаем ошибку
			return errors.Forbidden.New(ctx, "Access denied",
				errors.ParamsOption("UserID", userID, "AccountGroupID", accountGroupIDs),
				errors.HumanTextOption("Вы не имеете доступа к группе счетов с ID = %v", accountGroupID),
				errors.SkipPreviousCallerOption(),
			)
		}
	}

	return nil

}
