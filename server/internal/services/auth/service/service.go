package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"server/internal/services/transactor"
	userModel "server/internal/services/user/model"
	userRepository "server/internal/services/user/repository"
	userRepoModel "server/internal/services/user/repository/model"
)

var tracer = otel.Tracer("/server/internal/services/auth/service")

var _ UserRepository = new(userRepository.UserRepository)
var _ GeneralRepository = new(transactor.Transactor)

type UserRepository interface {
	GetUsers(context.Context, userModel.GetUsersReq) ([]userModel.User, error)
	CreateUser(context.Context, userModel.CreateReq) (uint32, error)

	CreateDevice(context.Context, userModel.Device) (uint32, error)
	DeleteDevice(ctx context.Context, userID uint32, deviceID string) error
	UpdateDevice(context.Context, userRepoModel.UpdateDeviceReq) error
	GetDevices(context.Context, userRepoModel.GetDevicesReq) ([]userModel.Device, error)
}

type GeneralRepository interface {
	WithinTransaction(ctx context.Context, callback func(ctx context.Context) error) error
}

type AuthService struct {
	userRepository    UserRepository
	generalRepository GeneralRepository
	generalSalt       []byte
}

func NewAuthService(
	userRepository UserRepository,
	generalRepository GeneralRepository,
	generalSalt []byte,

) *AuthService {
	return &AuthService{
		userRepository:    userRepository,
		generalRepository: generalRepository,
		generalSalt:       generalSalt,
	}
}
