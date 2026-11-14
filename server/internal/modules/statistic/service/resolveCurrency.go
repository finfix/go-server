package service

import (
	"context"

	"github.com/google/uuid"

	accountGroupModel "server/internal/modules/accountGroup/model"
	userModel "server/internal/modules/user/model"
	"server/internal/utils/necessary"
)

// resolveCurrency выбирает валюту ответа
func (s *StatisticService) resolveCurrency(ctx context.Context, userInformation necessary.NecessaryUserInformation, accountGroupIDs []uuid.UUID, currency string) (string, error) {

	// Явно переданная валюта
	if currency != "" {
		return currency, nil
	}

	// Валюта единственной выбранной группы
	if len(accountGroupIDs) == 1 {
		groups, err := s.accountGroupService.GetAccountGroups(ctx, accountGroupModel.GetAccountGroupsReq{
			Necessary:       userInformation,
			AccountGroupIDs: accountGroupIDs,
		})
		if err != nil {
			return "", err
		}
		if len(groups) == 1 {
			return groups[0].Currency, nil
		}
	}

	// Валюта пользователя по умолчанию
	user, err := s.userService.GetUser(ctx, userModel.GetUserReq{Necessary: userInformation}) //nolint:exhaustruct
	if err != nil {
		return "", err
	}
	return user.DefaultCurrency, nil
}
