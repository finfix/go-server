package service

import (
	"context"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
	"go.opentelemetry.io/otel"

	"pkg/errors"
	"pkg/log"

	"server/internal/services/pushNotificator/model"
)

var tracer = otel.Tracer("/server/internal/services/pushNotificator/service")

type PushNotificatorService struct {
	apns *apns2.Client
	isOn bool
}

func NewPushNotificatorService(isOn bool, apnsCredentials model.APNsCredentials) (*PushNotificatorService, error) {

	if !isOn {
		log.Warning(context.Background(), "SendNotification notificator is off")
		return &PushNotificatorService{
			isOn: isOn,
			apns: nil,
		}, nil
	}

	authKey, err := token.AuthKeyFromFile(apnsCredentials.KeyFilePath)
	if err != nil {
		return nil, errors.InternalServer.Wrap(err)
	}

	apnsClient := apns2.NewTokenClient(&token.Token{ // nolint:exhaustruct
		AuthKey: authKey,
		KeyID:   apnsCredentials.KeyID,
		TeamID:  apnsCredentials.TeamID,
	})
	apnsClient.Host = apns2.HostProduction

	return &PushNotificatorService{
		isOn: isOn,
		apns: apnsClient,
	}, nil
}
