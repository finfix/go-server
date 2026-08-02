package service

import (
	"context"

	"pkg/jwtManager"
	"server/internal/utils/auth"
	"server/internal/utils/errors"

	"github.com/google/uuid"

	"server/internal/modules/auth/model"
	"server/internal/modules/auth/service/utils"
	userRepoModel "server/internal/modules/user/repository/model"
)

// RefreshTokens обновляет токены доступа в базе данных
func (s *AuthService) RefreshTokens(ctx context.Context, req model.RefreshTokensReq) (newTokens model.RefreshTokensRes, err error) {
	ctx, span := tracer.Start(ctx, "RefreshTokens")
	defer span.End()

	// RefreshTokens не проходит через auth-интерцептор (access-токен на этот момент
	// может быть просрочен), поэтому UserID и DeviceID получаем не из контекста,
	// а из самого refresh-токена, переданного клиентом
	claims, err := jwtManager.ParseToken[auth.Claims](ctx, req.Token, jwtManager.RefreshToken)
	if err != nil {
		return newTokens, err
	}

	// Получаем девайс по идентификатору пользователя и девайса
	devices, err := s.userRepository.GetDevices(ctx, userRepoModel.GetDevicesReq{ // nolint:exhaustruct
		DeviceIDs: []string{claims.DeviceID},
		UserIDs:   []uuid.UUID{claims.UserID},
	})
	if err != nil {
		return newTokens, err
	}
	if len(devices) == 0 {
		return newTokens, errors.Unauthorized.New("Device not found").WithContextParams(ctx).
			WithCustomHumanText("Девайс не найден")
	}
	device := devices[0]

	// Сравниваем токен из базы данных с переданным пользователем токеном
	if req.Token != device.RefreshToken {
		return newTokens, errors.Forbidden.New("Auth is incorrect").
			WithContextParams(ctx)
	}

	// Создаем новую пару токенов
	newTokens.Tokens, err = utils.CreatePairTokens(ctx, claims.UserID, claims.DeviceID)
	if err != nil {
		return newTokens, err
	}

	// Обновляем данные у девайса
	if err = s.userRepository.UpdateDevice(ctx, userRepoModel.UpdateDeviceReq{
		UserID:       claims.UserID,
		DeviceID:     claims.DeviceID,
		RefreshToken: &newTokens.Tokens.RefreshToken,
		DeviceInformation: userRepoModel.UpdateDeviceInformationReq{
			VersionOS: &req.Device.VersionOS,
			UserAgent: &req.Device.UserAgent,
			IPAddress: &req.Device.IPAddress,
		},
		ApplicationInformation: userRepoModel.UpdateApplicationInformationReq{
			BundleID: &req.Application.BundleID,
			Version:  &req.Application.Version,
			Build:    &req.Application.Build,
		},
		NotificationToken: nil,
	}); err != nil {
		return newTokens, err
	}

	// Возвращаем пару токенов клиенту
	return newTokens, nil
}