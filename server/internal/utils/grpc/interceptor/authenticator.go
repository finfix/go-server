package interceptor

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc"

	"pkg/jwtManager"
	"server/internal/utils/auth"
	"server/internal/utils/contextKeys"
	"server/internal/utils/errors"
)

type AuthInterceptor struct {
	disableAuthorizationPaths []string
}

func NewAuthInterceptor(accessibleRoutes []string) *AuthInterceptor {
	return &AuthInterceptor{accessibleRoutes}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, err := interceptor.authorize(ctx, info.FullMethod, req)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		_, err := interceptor.authorize(stream.Context(), info.FullMethod, stream)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string, req any) (context.Context, error) {

	// Игнорируем авторизацию для заранее перечисленных методов
	for _, accessibleRoute := range interceptor.disableAuthorizationPaths {
		if method == accessibleRoute {
			// everyone can access
			return ctx, nil
		}
	}

	type authNecessaryFields struct {
		AccessToken string `json:"accessToken"`
	}

	// Сериализуем запрос в JSON
	requestJson, err := json.Marshal(req)
	if err != nil {
		return ctx, errors.BadRequest.Wrap(err)
	}

	// Десериализуем JSON в структуру с интересующим полем
	var authFields authNecessaryFields
	err = json.Unmarshal(requestJson, &authFields)
	if err != nil {
		return ctx, errors.BadRequest.Wrap(err)
	}

	// Парсим токен
	claims, err := jwtManager.ParseToken[auth.Claims](authFields.AccessToken, jwtManager.AccessToken)
	if err != nil {
		return ctx, errors.Unauthorized.Wrap(err)
	}

	// Добавляем в контекст UUID пользователя
	// Устанавливаем данные из токена в контекст
	ctx = contextKeys.SetUserID(ctx, claims.UserID)
	ctx = contextKeys.SetDeviceID(ctx, claims.DeviceID)

	return ctx, nil
}
