package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shopspring/decimal"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/sync/errgroup"

	"pkg/database/postgresql"
	"pkg/errors"
	"pkg/http/middleware"
	"pkg/http/router"
	"pkg/http/server"
	"pkg/jwtManager"
	"pkg/log"
	"pkg/log/model"
	"pkg/migrator"
	"pkg/panicRecover"
	"pkg/stackTrace"
	"pkg/trace"
	"server/internal/config"
	_ "server/internal/docs"
	accountEndpoint "server/internal/services/account/endpoint"
	accountRepository "server/internal/services/account/repository"
	accountService "server/internal/services/account/service"
	accountGroupEndpoint "server/internal/services/accountGroup/endpoint"
	accountGroupRepository "server/internal/services/accountGroup/repository"
	accountGroupService "server/internal/services/accountGroup/service"
	accountPermisssionsRepository "server/internal/services/accountPermissions/repository"
	accountPermisssionsService "server/internal/services/accountPermissions/service"
	authEndpoint "server/internal/services/auth/endpoint"
	authService "server/internal/services/auth/service"
	pushNotificatorModel "server/internal/services/pushNotificator/model"
	pushNotificatorService "server/internal/services/pushNotificator/service"
	"server/internal/services/scheduler"
	settingsEndpoint "server/internal/services/settings/endpoint"
	settingsRepository "server/internal/services/settings/repository"
	settingsService "server/internal/services/settings/service"
	tagEndpoint "server/internal/services/tag/endpoint"
	tagRepository "server/internal/services/tag/repository"
	tagService "server/internal/services/tag/service"
	"server/internal/services/tgBot/service"
	transactionEndpoint "server/internal/services/transaction/endpoint"
	transactionRepository "server/internal/services/transaction/repository"
	transactionService "server/internal/services/transaction/service"
	"server/internal/services/transactor"
	userEndpoint "server/internal/services/user/endpoint"
	userRepository "server/internal/services/user/repository"
	userService "server/internal/services/user/service"
	"server/migrations"
)

// @title COIN Server Documentation
// @version @{version} (build @{build})
// @description API Documentation for Coin
// @contact.name Ilia Ivanov
// @contact.email bonavii@icloud.com
// @contact.url

// @securityDefinitions.apikey AuthJWT
// @in header
// @name Authorization
// @description JWT-токен авторизации

//go:generate go install github.com/swaggo/swag/cmd/swag@v1.8.2
//go:generate go mod download
//go:generate swag init -o docs --parseInternal --parseDependency

const version = "@{version}"
const build = "@{build}"

func main() {
	if err := run(); err != nil {
		log.Fatal(context.Background(), err)
	}
}

func run() error {

	// Основной контекст приложения
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Перехватываем возможную панику
	defer panicRecover.PanicRecover(func(err error) {
		log.Fatal(ctx, err)
	})

	// Парсим флаги
	logFormat := flag.String("log-format", string(log.JSONFormat), "text - Human readable string\njson - JSON format")
	envMode := flag.String("env-mode", "local", "Environment mode for log label: test, prod")
	flag.Parse()

	var logHandlers []log.Handler
	switch *logFormat {
	case "text":
		logHandlers = append(logHandlers, log.NewConsoleHandler(os.Stdout, log.LevelDebug))
	case "json":
		logHandlers = append(logHandlers, log.NewJSONHandler(os.Stdout, log.LevelDebug))
	}

	// Получаем имя хоста
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	// Инициализируем логгер
	log.Init(
		model.SystemInfo{
			Hostname: hostname,
			Version:  version,
			Build:    build,
			Env:      *envMode,
		},
		logHandlers...,
	)

	// Получаем конфиг
	log.Info(ctx, "Получаем конфиг")
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	// Инициализируем все синглтоны
	log.Info(ctx, "Инициализируем синглтоны")
	if err = initSingletons(cfg); err != nil {
		return err
	}

	log.Info(ctx, "Инициализируем трейсер")
	if err = trace.StartTracing(cfg.Tracer, cfg.ServiceName); err != nil {
		return err
	}

	// Подключаемся к базе данных
	log.Info(ctx, "Подключаемся к БД")
	postrgreSQL, err := postgresql.NewClientSQL(cfg.Repository, cfg.DBName)
	if err != nil {
		return err
	}
	defer postrgreSQL.Close()

	// Запускаем миграции в базе данных
	// TODO: Подумать, как откатывать миграции при ошибках
	log.Info(ctx, "Запускаем миграции")
	postgreSQLMigrator := migrator.NewMigrator(
		postrgreSQL,
		migrator.MigratorConfig{
			EmbedMigrations: migrations.EmbedMigrationsPostgreSQL,
			Dir:             "pgsql",
		},
	)
	if err = postgreSQLMigrator.Up(ctx); err != nil {
		return err
	}

	// Регистрируем репозитории
	transactor := transactor.NewTransactor(postrgreSQL)
	accountGroupRepository := accountGroupRepository.NewAccountGroupRepository(postrgreSQL)
	accountRepository := accountRepository.NewAccountRepository(postrgreSQL)
	tagRepository := tagRepository.NewTagRepository(postrgreSQL)
	transactionRepository := transactionRepository.NewTransactionRepository(postrgreSQL)
	settingsRepository := settingsRepository.NewSettingsRepository(postrgreSQL)
	userRepository := userRepository.NewUserRepository(postrgreSQL)
	accountPermissionsRepository := accountPermisssionsRepository.NewAccountPermissionsRepository(postrgreSQL)

	// Регистрируем сервисы
	log.Info(ctx, "Инициализируем Telegram-бота")
	tgBotService, err := service.NewTgBotService(cfg.Telegram.Token, cfg.Telegram.ChatID, cfg.Telegram.Enabled)
	if err != nil {
		return err
	}
	if cfg.Telegram.Enabled {
		defer tgBotService.Bot.Close()
	}

	log.Info(ctx, "Инициализируем сервис пушей")
	pushNotificatorService, err := pushNotificatorService.NewPushNotificatorService(cfg.Notifications.Enabled, pushNotificatorModel.APNsCredentials{
		TeamID:      cfg.Notifications.APNs.TeamID,
		KeyID:       cfg.Notifications.APNs.KeyID,
		KeyFilePath: cfg.Notifications.APNs.KeyFilePath,
	})
	if err != nil {
		return err
	}

	accountPermissionsService := accountPermisssionsService.NewAccountPermissionsService(accountPermissionsRepository)

	userService := userService.NewUserService(
		userRepository,
		transactor,
		pushNotificatorService,
		[]byte(cfg.GeneralSalt),
	)

	accountGroupService := accountGroupService.NewAccountGroupService(
		accountGroupRepository,
		transactor,
		userService,
	)

	accountService := accountService.NewAccountService(
		accountRepository,
		transactor,
		transactionRepository,
		userRepository,
		accountPermissionsService,
		accountGroupService,
		userService,
	)

	tagService := tagService.NewTagService(
		tagRepository,
		transactor,
		userService,
		accountGroupService,
	)

	transactionService := transactionService.NewTransactionService(
		transactionRepository,
		accountRepository,
		transactor,
		accountPermissionsService,
		tagRepository,
		userService,
		accountService,
		tagService,
	)

	settingsService := settingsService.NewSettingsService(
		settingsRepository,
		userService,
		tgBotService,
		settingsService.Version{
			Version: version,
			Build:   build,
		},
		settingsService.Credentials{
			CurrencyProviderAPIKey: cfg.APIKeys.CurrencyProvider,
		},
	)

	authService := authService.NewAuthService(
		userRepository,
		transactor,
		[]byte(cfg.GeneralSalt),
	)

	log.Info(ctx, "Запускаем планировщик")
	if err = scheduler.NewScheduler(settingsService).Start(); err != nil {
		return err
	}

	r := router.NewRouter()
	accountEndpoint.MountAccountEndpoints(r, accountService)                // ANY /account*
	accountGroupEndpoint.MountAccountGroupEndpoints(r, accountGroupService) // ANY /accountGroup*
	transactionEndpoint.MountTransactionEndpoints(r, transactionService)    // ANY /transaction*
	tagEndpoint.MountTagEndpoints(r, tagService)                            // ANY /tag*
	authEndpoint.MountAuthEndpoints(r, authService)                         // ANY /auth*
	userEndpoint.MountUserEndpoints(r, userService)                         // ANY /user*
	settingsEndpoint.MountSettingsEndpoints(r, settingsService)             // ANY /settings*
	r.Mount("/swagger", httpSwagger.WrapHandler)

	server, err := server.GetDefaultServer(cfg.HTTP, r)
	if err != nil {
		return err
	}

	// Создаем wait группу
	eg, ctx := errgroup.WithContext(ctx)

	// Запускаем HTTP-сервер
	eg.Go(func() error { return server.Serve(ctx) })

	// Запускаем горутину, ожидающую завершение контекста
	eg.Go(func() error {

		// Если контекст завершился, значит процесс убили
		<-ctx.Done()

		// Плавно завершаем работу сервера
		server.Shutdown(ctx)

		return nil
	})

	// Ждем завершения контекста или ошибок в горутинах
	return eg.Wait()
}

func initSingletons(cfg config.Config) error {

	stackTrace.Init(cfg.ServiceName)

	// Конфигурируем decimal, чтобы в JSON не было кавычек
	decimal.MarshalJSONWithoutQuotes = true

	// Инициализируем JWT-менеджер
	accessTokenTTL, err := time.ParseDuration(cfg.Token.AccessTokenTTL)
	if err != nil {
		return errors.InternalServer.Wrap(err)
	}
	refreshTokenTTL, err := time.ParseDuration(cfg.Token.RefreshTokenTTL)
	if err != nil {
		return errors.InternalServer.Wrap(err)
	}
	jwtManager.Init([]byte(cfg.Token.SigningKey), accessTokenTTL, refreshTokenTTL)

	if err = middleware.Init(cfg.ServiceName); err != nil {
		return err
	}

	return nil
}
