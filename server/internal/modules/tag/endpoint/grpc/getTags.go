package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/tag/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// GetTags получение всех тегов
func (e *TagEndpoint) GetTags(ctx context.Context, r *proto.GetTagsRequest) (*proto.GetTagsResponse, error) {
	res := &proto.GetTagsResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoGetTagsReq{GetTagsRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Parse necessary information from context
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	// Validate request
	if err := validator.Validate(req); err != nil {
		return res, err
	}

	// Call service method
	tags, err := e.tagService.GetTags(ctx, req)
	if err != nil {
		return res, err
	}

	// Convert tags to proto format
	protoTags := make([]*proto.Tag, 0, len(tags))
	for _, tag := range tags {
		protoTag, err := tag.ConvertToProto()
		if err != nil {
			return res, err
		}
		protoTags = append(protoTags, protoTag)
	}

	res.Tags = protoTags
	return res, nil
}
