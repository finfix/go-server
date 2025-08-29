package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/tag/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// DeleteTag удаление тега
func (e *TagEndpoint) DeleteTag(ctx context.Context, r *proto.DeleteTagRequest) (*proto.DeleteTagResponse, error) {
	res := &proto.DeleteTagResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoDeleteTagReq{DeleteTagRequest: r}.ConvertToModel()
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
	err = e.tagService.DeleteTag(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}
