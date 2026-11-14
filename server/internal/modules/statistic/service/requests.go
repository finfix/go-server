package service

import (
	"server/internal/modules/statistic/model"
	tagModel "server/internal/modules/tag/model"
)

// tagGetReq — запрос тегов выбранных групп
func tagGetReq(req model.GetStatisticReq, data dataset) tagModel.GetTagsReq {
	return tagModel.GetTagsReq{Necessary: req.Necessary, AccountGroupIDs: data.accountGroupIDs} //nolint:exhaustruct
}
