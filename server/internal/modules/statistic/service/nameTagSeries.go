package service

import (
	"context"

	"server/internal/modules/statistic/model"
)

// nameTagSeries подписывает ряды разбивки по тегам их названиями
func (s *StatisticService) nameTagSeries(ctx context.Context, req model.GetStatisticReq, data dataset, series map[string]*model.Series) error {
	tags, err := s.tagService.GetTags(ctx, tagGetReq(req, data))
	if err != nil {
		return err
	}
	for _, tag := range tags {
		if item, ok := series[tag.ID.String()]; ok {
			item.Name = tag.Name
		}
	}
	if item, ok := series[keyUntagged]; ok {
		item.Name = "Без тега"
	}
	return nil
}
