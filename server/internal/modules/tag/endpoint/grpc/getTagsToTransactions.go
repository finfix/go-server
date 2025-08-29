package grpc

import (
	"context"
	"server/internal/utils/necessary"

	"pkg/validator"
	"server/internal/modules/tag/model"

	proto "github.com/finfix/go-server-grpc/proto"
)

// GetTagsToTransactions получение связей тегов с транзакциями
func (e *TagEndpoint) GetTagsToTransactions(ctx context.Context, r *proto.GetTagsToTransactionsRequest) (*proto.GetTagsToTransactionsResponse, error) {
	res := &proto.GetTagsToTransactionsResponse{}

	// Convert proto request to internal model
	req, err := model.ProtoGetTagsToTransactionsReq{GetTagsToTransactionsRequest: r}.ConvertToModel()
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
	tagsToTransactions, err := e.tagService.GetTagsToTransactions(ctx, req)
	if err != nil {
		return res, err
	}

	// Convert to proto format
	protoTagsToTransactions := make([]*proto.TagToTransaction, 0, len(tagsToTransactions))
	for _, tagToTransaction := range tagsToTransactions {
		protoTagToTransaction, err := tagToTransaction.ConvertToProto()
		if err != nil {
			return res, err
		}
		protoTagsToTransactions = append(protoTagsToTransactions, protoTagToTransaction)
	}

	res.TagsToTransactions = protoTagsToTransactions
	return res, nil
}
