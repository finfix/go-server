package service

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	pushNotificatorModel "server/internal/modules/pushNotificator/model"
	pushNotificatorService "server/internal/modules/pushNotificator/service"
	"server/internal/modules/transactor"
	userModel "server/internal/modules/user/model"
	userRepository "server/internal/modules/user/repository"
	userRepoModel "server/internal/modules/user/repository/model"
)

var tracer = otel.Tracer("/server/internal/modules/user/service")

var _ UserRepository = new(userRepository.UserRepository)
var _ GeneralRepository = new(transactor.Transactor)
var _ PushNotificatorService = new(pushNotificatorService.PushNotificatorService)

type GeneralRepository interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
}

type UserRepository interface {
	CreateUser(context.Context, userModel.CreateReq) (uuid.UUID, error)
	GetUsers(context.Context, userModel.GetUsersReq) ([]userModel.User, error)
	UpdateUser(context.Context, userRepoModel.UpdateUserReq) error

	LinkUserToAccountGroup(context.Context, uuid.UUID, uuid.UUID) error

	GetDevices(context.Context, userRepoModel.GetDevicesReq) ([]userModel.Device, error)
	UpdateDevice(context.Context, userRepoModel.UpdateDeviceReq) error

	GetAccessedAccountGroups(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type PushNotificatorService interface {
	SendNotification(ctx context.Context, req pushNotificatorModel.SendNotificationReq) (string, error)
}

type UserService struct {
	userRepository    UserRepository
	generalRepository GeneralRepository
	pushNotificator   PushNotificatorService
	generalSalt       []byte
}

func NewUserService(
	userRepository UserRepository,
	generalRepository GeneralRepository,
	pushNotificator PushNotificatorService,
	generalSalt []byte,
) *UserService {
	return &UserService{
		userRepository:    userRepository,
		generalRepository: generalRepository,
		pushNotificator:   pushNotificator,
		generalSalt:       generalSalt,
	}
}
