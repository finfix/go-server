package model

import (
	"context"

	"server/internal/utils/errors"
)

// GroupBy — разбивка рядов доходов и расходов
type GroupBy string

const (
	GroupByAccount GroupBy = "account" // По счетам (категориям)
	GroupByTag     GroupBy = "tag"     // По тегам (подкатегориям)
)

// Validate проверяет, что разбивка известна
func (g GroupBy) Validate(ctx context.Context) error {
	switch g {
	case GroupByAccount, GroupByTag:
		return nil
	}
	return errors.BadRequest.New("unknown group_by").WithParams("groupBy", string(g)).WithContextParams(ctx)
}
