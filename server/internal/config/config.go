package config

import (
	"pkg/env"
	"pkg/log"
	"pkg/trace"

	"pkg/database/pgsql"
)

// Config - общая структура конфига
type Config struct {

	// Адрес для http-сервера
	Port struct {
		HTTP string `env:"PORT_HTTP"`
		GRPC string `env:"PORT_GRPC"`
	}

	// Данные базы данных
	Pgsql pgsql.PgsqlConfigEnv

	Tracer trace.TracerConfig

	// Информация для JWT-токенов
	Auth struct {
		AccessTokenTTL  string `env:"AUTH_ACCESS_TOKEN_TTL"`
		RefreshTokenTTL string `env:"AUTH_REFRESH_TOKEN_TTL"`
		SigningKey      string `env:"AUTH_TOKEN_SIGNING_KEY"`
		// Информация для шифрования паролей
		GeneralSalt string `env:"SHA_SALT"`
	}

	// Ключи для работы с внешним API
	APIKeys struct {
		CurrencyProvider string `env:"API_KEY_CURRENCY_PROVIDER"`
	}

	// Доступы к телеграм-боту
	Telegram struct {
		IsEnabled bool   `env:"TG_BOT_ENABLED"`
		Token     string `env:"TG_BOT_TOKEN"`
		ChatID    int64  `env:"TG_CHAT_ID"`
	}

	Notifications struct {
		Enabled bool `env:"NOTIFICATIONS_ENABLED"`
		APNs    struct {
			TeamID      string `env:"NOTIFICATIONS_APNS_TEAM_ID"`
			KeyID       string `env:"NOTIFICATIONS_APNS_KEY_ID"`
			KeyFilePath string `env:"NOTIFICATIONS_APNS_KEY_FILE_PATH"`
		}
	}

	ServiceName string `env:"SERVICE_NAME"`
	Environment string `env:"ENVIRONMENT"`

	Logger log.LoggerSettingsEnv

	StackTraceEnabled bool `env:"STACK_TRACE_ENABLED"`
}

func Load() Config {
	return env.Load[Config]()
}
