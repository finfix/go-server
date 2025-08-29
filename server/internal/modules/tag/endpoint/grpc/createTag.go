package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/tag/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// CreateTag создание тега
func (e *TagEndpoint) CreateTag(ctx context.Context, r *proto.CreateTagRequest) (*proto.CreateTagResponse, error) {
	res := &proto.CreateTagResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoCreateTagReq{CreateTagRequest: r}.ConvertToModel()
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
	createRes, err := e.tagService.CreateTag(ctx, req)
	if err != nil {
		return res, err
	}

	// Convert response to proto and return
	return createRes.ConvertToProto(), nil
}
