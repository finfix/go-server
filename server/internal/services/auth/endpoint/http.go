package endpoint

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"

	"pkg/errors"
	"pkg/http/chain"
	"pkg/jwtManager"
	"server/internal/services/auth/model"
	"server/internal/utils/auth"
	"server/internal/utils/contextKeys"
)

var tracer = otel.Tracer("/server/internal/services/auth/endpoint")

type endpoint struct {
	service authService
}

type authService interface {
	SignIn(context.Context, model.SignInReq) (model.AuthRes, error)
	SignUp(context.Context, model.SignUpReq) (model.AuthRes, error)
	SignOut(context.Context, model.SignOutReq) error
	RefreshTokens(context.Context, model.RefreshTokensReq) (model.RefreshTokensRes, error)
}

func MountAuthEndpoints(mux *chi.Mux, service authService) {
	mux.Mount("/auth", newAuthEndpoint(service))
}

func newAuthEndpoint(service authService) http.Handler {

	s := &endpoint{
		service: service,
	}

	options := []chain.Option{
		chain.Before(deviceIDValidator),
	}

	r := chi.NewRouter()

	r.Method(http.MethodPost, "/signIn", chain.NewChain(s.signIn, options...))
	r.Method(http.MethodPost, "/signUp", chain.NewChain(s.signUp, options...))
	r.Method(http.MethodPost, "/signOut", chain.NewChain(s.signOut, options...))
	r.Method(http.MethodPost, "/refreshTokens", chain.NewChain(s.refreshTokens, append(options, chain.Before(extractDataFromToken))...))
	return r
}

func deviceIDValidator(ctx context.Context, r *http.Request) (context.Context, error) {
	_, span := tracer.Start(ctx, "DefaultDeviceIDValidator")
	defer span.End()

	// Получаем DeviceID из заголовка
	deviceID := r.Header.Get("DeviceID")
	if deviceID == "" {
		return ctx, errors.BadRequest.New(ctx, "DeviceID is empty")
	}

	// Сохраняем DeviceID в контекст
	ctx = contextKeys.SetDeviceID(ctx, deviceID)

	return ctx, nil
}

func extractDataFromToken(ctx context.Context, r *http.Request) (context.Context, error) {

	// Проводим авторизацию
	ctx, err := auth.DefaultAuthorization(ctx, r)
	if err != nil {

		// Если ошибка истекшего токена, то это ок, так как мы смогли его распарсить и получить оттуда данные
		if errors.Is(err, jwtManager.ErrTokenExpired) {
			return ctx, nil
		}

		return ctx, err
	}

	return ctx, nil
}
