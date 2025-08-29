package grpc

import (
	"context"

	"server/internal/modules/settings/model"
	settingsService "server/internal/modules/settings/service"

	"github.com/finfix/go-server-grpc/proto"
)

var _ SettingsService = new(settingsService.SettingsService)

type SettingsService interface {
	GetCurrencies(context.Context) (model.GetCurrenciesRes, error)
	GetIcons(context.Context) (model.GetIconsRes, error)
	UpdateCurrencies(context.Context, model.UpdateCurrenciesReq) error
	SendNotification(context.Context, model.SendNotificationReq) (model.SendNotificationRes, error)
	GetVersion(context.Context, model.GetVersionReq) (model.GetVersionRes, error)
}

var _ proto.SettingsEndpointServer = new(SettingsEndpoint)

type SettingsEndpoint struct {
	proto.UnsafeSettingsEndpointServer
	settingsService SettingsService
}

func (e *SettingsEndpoint) UpdateCurrencies(ctx context.Context, request *proto.UpdateCurrenciesRequest) (*proto.UpdateCurrenciesResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (e *SettingsEndpoint) SendNotification(ctx context.Context, request *proto.SendNotificationRequest) (*proto.SendNotificationResponse, error) {
	// TODO implement me
	panic("implement me")
}

func NewSettingsEndpoint(settingsService SettingsService) *SettingsEndpoint {
	return &SettingsEndpoint{
		settingsService: settingsService,
	}
}
