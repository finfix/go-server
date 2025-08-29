package grpc

import (
	"context"

	proto "github.com/finfix/go-server-grpc/proto"
)

// GetIcons получение текущей версии сервера
func (e *SettingsEndpoint) GetIcons(ctx context.Context, r *proto.GetIconsRequest) (*proto.GetIconsResponse, error) {
	res := &proto.GetIconsResponse{}

	// Call service method
	icons, err := e.settingsService.GetIcons(ctx)
	if err != nil {
		return res, err
	}

	return icons.ConvertToProto()
}
