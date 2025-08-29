package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"server/internal/modules/accountPermissions/model"
)

var tracer = otel.Tracer("/server/internal/modules/accountPermissions/service")

type AccountPermissionsService struct {
	accountPermissionsRepository AccountPermissionsRepository
}

type AccountPermissionsRepository interface {
	GetAccountPermissions(context.Context) (model.PermissionSet, error)
}

func NewAccountPermissionsService(accountPermissionsRepository AccountPermissionsRepository) *AccountPermissionsService {
	return &AccountPermissionsService{
		accountPermissionsRepository: accountPermissionsRepository,
	}
}
