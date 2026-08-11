package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	auditLogModel "server/internal/modules/auditLog/model"
	auditLogService "server/internal/modules/auditLog/service"
	settingsModel "server/internal/modules/settings/model"
	settingsRepository "server/internal/modules/settings/repository"
	tgBotModel "server/internal/modules/tgBot/model"
	tgBotService "server/internal/modules/tgBot/service"
	"server/internal/modules/transactor"
	userModel "server/internal/modules/user/model"
	userService "server/internal/modules/user/service"
)

var tracer = otel.Tracer("/server/internal/modules/settings/service")

var _ SettingsRepository = new(settingsRepository.SettingsRepository)

type SettingsRepository interface {
	UpdateCurrencies(ctx context.Context, rates map[string]decimal.Decimal) error
	GetCurrencies(context.Context) ([]settingsModel.Currency, error)
	GetIcons(context.Context) ([]settingsModel.Icon, error)
	GetVersion(context.Context, settingsModel.GetVersionReq) (settingsModel.Version, error)
}

var _ UserService = new(userService.UserService)

type UserService interface {
	SendNotification(ctx context.Context, userID uuid.UUID, push userModel.Notification) (uint8, error)
	GetUsers(ctx context.Context, filters userModel.GetUsersReq) (users []userModel.User, err error)
}

var _ TgBotService = new(tgBotService.TgBotService)

type TgBotService interface {
	SendMessage(context.Context, tgBotModel.SendMessageReq) error
}

var _ Transactor = new(transactor.Transactor)

type Transactor interface {
	WithinTransaction(ctx context.Context, callback func(context.Context) error) error
}

var _ AuditLogService = new(auditLogService.AuditLogService)

// AuditLogService - интерфейс сервиса аудит-лога
type AuditLogService interface {
	TrackMutation(context.Context, auditLogModel.TrackMutationReq) (uint32, error)
}

type Credentials struct {
	CurrencyProviderAPIKey string
}

type Version struct {
	Version string
	Build   string
}

type SettingsService struct {
	settingsRepository SettingsRepository
	userService        UserService
	tgBot              TgBotService
	transactor         Transactor
	auditLogService    AuditLogService
	credentials        Credentials
	version            Version
}

func NewSettingsService(
	settingsRepository SettingsRepository,
	userService UserService,
	tgBot TgBotService,
	transactor Transactor,
	auditLogService AuditLogService,
	version Version,
	credentials Credentials,
) *SettingsService {
	return &SettingsService{
		settingsRepository: settingsRepository,
		userService:        userService,
		tgBot:              tgBot,
		transactor:         transactor,
		auditLogService:    auditLogService,
		credentials:        credentials,
		version:            version,
	}
}
