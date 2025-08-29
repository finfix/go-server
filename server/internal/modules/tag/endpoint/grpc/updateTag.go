package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/tag/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// UpdateTag обновление тега
func (e *TagEndpoint) UpdateTag(ctx context.Context, r *proto.UpdateTagRequest) (*proto.UpdateTagResponse, error) {
	res := &proto.UpdateTagResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoUpdateTagReq{UpdateTagRequest: r}.ConvertToModel()
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
	err = e.tagService.UpdateTag(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}
