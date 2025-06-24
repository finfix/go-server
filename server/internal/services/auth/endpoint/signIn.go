package endpoint

import (
	"context"
	"net/http"

	"pkg/http/decoder"
	"pkg/validator"
	"server/internal/services/auth/model"
	"server/internal/utils/contextKeys"
	"server/internal/utils/errors"
)

// @Summary Авторизация пользователя по логину и паролю
// @Tags auth
// @Accept json
// @Param Body body model.SignInReq true "model.SignInReq"
// @Param DeviceID header string true "Нужен для идентификации устройства"
// @Produce json
// @Success 200 {object} model.AuthRes
// @Failure 400,404,500 {object} errors.Error
// @Router /auth/signIn [post]
func (s *endpoint) signIn(ctx context.Context, r *http.Request) (any, error) {
	ctx, span := tracer.Start(ctx, "signIn")
	defer span.End()

	var req model.SignInReq

	deviceID := contextKeys.GetDeviceID(ctx)
	if deviceID == nil {
		return nil, errors.BadRequest.New("DeviceID не задан").WithContextParams(ctx)
	}
	req.DeviceID = *deviceID

	// Декодируем запрос
	if err := decoder.Decode(ctx, r, &req, decoder.DecodeJSON); err != nil {
		return nil, err
	}

	req.Device.IPAddress = r.Header.Get("X-Real-IP")
	req.Device.UserAgent = r.Header.Get("User-Agent")

	// Валидируем запрос
	if err := validator.Validate(req); err != nil {
		return nil, err
	}

	// Вызываем метод сервиса
	return s.service.SignIn(ctx, req)
}
