package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"

	"server/internal/modules/transactor"
	userModel "server/internal/modules/user/model"
	userRepository "server/internal/modules/user/repository"
	userRepoModel "server/internal/modules/user/repository/model"
)

var tracer = otel.Tracer("/server/internal/modules/auth/service")

var _ UserRepository = new(userRepository.UserRepository)
var _ GeneralRepository = new(transactor.Transactor)

type UserRepository interface {
	GetUsers(context.Context, userModel.GetUsersReq) ([]userModel.User, error)
	CreateUser(context.Context, userModel.CreateReq) error

	CreateDevice(context.Context, userModel.Device) error
	DeleteDevice(ctx context.Context, userID uuid.UUID, deviceID string) error
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
