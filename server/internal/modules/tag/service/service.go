package service

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	accountGroupService "server/internal/modules/accountGroup/service"
	tagModel "server/internal/modules/tag/model"
	tagRepository "server/internal/modules/tag/repository"
	tagRepoModel "server/internal/modules/tag/repository/model"
	"server/internal/modules/transactor"
	userService "server/internal/modules/user/service"
)

var tracer = otel.Tracer("/server/internal/modules/tag/service")

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
	GetAccessedAccountGroups(ctx context.Context, userID uuid.UUID) (accesses []uuid.UUID, err error)
}

var _ Transactor = &transactor.Transactor{}

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
}

var _ TagRepository = &tagRepository.TagRepository{}

type TagRepository interface {
	CreateTag(context.Context, tagRepoModel.CreateTagReq) (uuid.UUID, error)
	UpdateTag(context.Context, tagModel.UpdateTagReq) error
	DeleteTag(ctx context.Context, id, userID uuid.UUID) error
	GetTags(context.Context, tagModel.GetTagsReq) (res []tagModel.Tag, err error)

	GetTagsToTransactions(ctx context.Context, req tagModel.GetTagsToTransactionsReq) ([]tagModel.TagToTransaction, error)

	CheckAccess(ctx context.Context, accountGroupIDs, tagIDs []uuid.UUID) error
}

var _ AccountGroupService = &accountGroupService.AccountGroupService{}

type AccountGroupService interface {
	CheckAccess(ctx context.Context, userID uuid.UUID, accountGroupIDs []uuid.UUID) error
}
