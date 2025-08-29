package grpc

import (
	"context"

	"server/internal/modules/tag/model"
	tagService "server/internal/modules/tag/service"

	proto "github.com/finfix/go-server-grpc/proto"
)

var _ TagService = new(tagService.TagService)

type TagService interface {
	CreateTag(context.Context, model.CreateTagReq) (model.CreateTagRes, error)
	GetTags(context.Context, model.GetTagsReq) ([]model.Tag, error)
	UpdateTag(context.Context, model.UpdateTagReq) error
	DeleteTag(context.Context, model.DeleteTagReq) error
	GetTagsToTransactions(context.Context, model.GetTagsToTransactionsReq) ([]model.TagToTransaction, error)
}

var _ proto.TagEndpointServer = new(TagEndpoint)

type TagEndpoint struct {
	proto.UnsafeTagEndpointServer
	tagService TagService
}

func NewTagEndpoint(tagService TagService) *TagEndpoint {
	return &TagEndpoint{
		tagService: tagService,
	}
}
