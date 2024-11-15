package service

import (
	"context"

	"go.opentelemetry.io/otel"

	accountGroupService "server/internal/services/accountGroup/service"
	tagModel "server/internal/services/tag/model"
	tagRepository "server/internal/services/tag/repository"
	tagRepoModel "server/internal/services/tag/repository/model"
	"server/internal/services/transactor"
	userService "server/internal/services/user/service"
)

var tracer = otel.Tracer("/server/internal/services/tag/service")

type TagService struct {
	tagRepository       TagRepository
	generalRepository   Transactor
	userService         UserService
	accountGroupService AccountGroupService
}

func NewTagService(
	tagRepository TagRepository,
	generalRepository Transactor,
	userService UserService,
	accountGroupService AccountGroupService,
) *TagService {
	return &TagService{
		tagRepository:       tagRepository,
		generalRepository:   generalRepository,
		userService:         userService,
		accountGroupService: accountGroupService,
	}
}

var _ UserService = &userService.UserService{}

type UserService interface {
	GetAccessedAccountGroups(ctx context.Context, userID uint32) (accesses []uint32, err error)
}

var _ Transactor = &transactor.Transactor{}

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
}

var _ TagRepository = &tagRepository.TagRepository{}

type TagRepository interface {
	CreateTag(context.Context, tagRepoModel.CreateTagReq) (uint32, error)
	UpdateTag(context.Context, tagModel.UpdateTagReq) error
	DeleteTag(ctx context.Context, id, userID uint32) error
	GetTags(context.Context, tagModel.GetTagsReq) (res []tagModel.Tag, err error)

	GetTagsToTransactions(ctx context.Context, req tagModel.GetTagsToTransactionsReq) ([]tagModel.TagToTransaction, error)

	CheckAccess(ctx context.Context, accountGroupIDs, tagIDs []uint32) error
}

var _ AccountGroupService = &accountGroupService.AccountGroupService{}

type AccountGroupService interface {
	CheckAccess(ctx context.Context, userID uint32, accountGroupIDs []uint32) error
}
