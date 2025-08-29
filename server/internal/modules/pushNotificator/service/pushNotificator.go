package service

import (
	"go.opentelemetry.io/otel"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"

	"pkg/log"
	"server/internal/utils/errors"

	"server/internal/modules/pushNotificator/model"
)

var tracer = otel.Tracer("/server/internal/modules/pushNotificator/service")

type PushNotificatorService struct {
	apns *apns2.Client
	isOn bool
}

func NewPushNotificatorService(isOn bool, apnsCredentials model.APNsCredentials) (*PushNotificatorService, error) {

	if !isOn {
		log.Warning("SendNotification notificator is off")
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
