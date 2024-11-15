package config

import (
	"github.com/caarlos0/env/v11"

	"pkg/database/postgresql"
	"pkg/errors"
	"pkg/trace"
)

// Config - общая структура конфига
type Config struct {

	// Адрес для http-сервера
	HTTP string `env:"LISTEN_HTTP" envDefault:":8080"`

	// Данные базы данных
	Repository postgresql.PostgreSQLConfig
	DBName     string `env:"DB_NAME" envDefault:"coin"`

	Tracer trace.TracerConfig

	// Информация для JWT-токенов
	Token struct {
		AccessTokenTTL  string `env:"AUTH_ACCESS_TOKEN_TTL" envDefault:"15m"`
		RefreshTokenTTL string `env:"AUTH_REFRESH_TOKEN_TTL" envDefault:"720h"`
		SigningKey      string `env:"AUTH_TOKEN_SIGNING_KEY" envDefault:"secret"`
	}

	// Информация для шифрования паролей
	GeneralSalt string `env:"SHA_SALT" envDefault:"secret"`

	// Ключи для работы с внешним API
	APIKeys struct {
		CurrencyProvider string `env:"API_KEY_CURRENCY_PROVIDER"  envDefault:"secret"`
	}

	// Доступы к телеграм-боту
	Telegram struct {
		Enabled bool   `env:"TG_BOT_ENABLED" envDefault:"false"`
		Token   string `env:"TG_BOT_TOKEN" envDefault:"secret"`
		ChatID  int64  `env:"TG_CHAT_ID" envDefault:"secret"`
	}

	Notifications struct {
		Enabled bool `env:"NOTIFICATIONS_ENABLED" envDefault:"false"`
		APNs    struct {
			TeamID      string `env:"NOTIFICATIONS_APNS_TEAM_ID" envDefault:"secret"`
			KeyID       string `env:"NOTIFICATIONS_APNS_KEY_ID" envDefault:"secret"`
			KeyFilePath string `env:"NOTIFICATIONS_APNS_KEY_FILE_PATH" envDefault:"secret"`
		}
	}

	ServiceName string `env:"SERVICE_NAME" envDefault:"coin"`
}

// GetConfig возвращает конфигурацию из .env файла
func GetConfig() (config Config, err error) {
	if err = env.Parse(&config); err != nil {
		return config, errors.InternalServer.Wrap(err)
	}
	return config, nil
}
