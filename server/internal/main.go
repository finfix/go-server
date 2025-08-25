package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pkg/log/model"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/pressly/goose/v3"
	"github.com/shopspring/decimal"
	httpSwagger "github.com/swaggo/http-swagger"

	"pkg/database/pgsql"
	"pkg/http/router"
	"pkg/http/server"
	"pkg/jwtManager"
	"pkg/log"
	"pkg/migrator"
	"pkg/panicRecover"
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
	"server/internal/utils/errors"
	pgsqlMigrations "server/migrations/pgsql"
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
const buildDate = "@{buildDate}"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// Основной контекст приложения
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Перехватываем возможную панику
	defer func() {
		panicRecover.PanicRecover(func(err error) {
			log.Fatal(err)
		})
	}()

	// Получаем конфиг
	conf := config.Load()

	// Получаем имя хоста
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	// Инициализируем логгер
	if err := log.InitDefaultLogger(
		model.SystemInfo{
			Version:     version,
			Build:       build,
			ServiceName: conf.ServiceName,
			Env:         conf.Environment,
			Hostname:    hostname,
			BuildDate:   buildDate,
		},
		conf.Logger,
	); err != nil {
		return err
	}

	// Инициализируем все синглтоны
	log.Info("Инициализируем синглтоны")
	if err = initSingletons(conf); err != nil {
		return err
	}

	log.Info("Инициализируем трейсер")
	if err = trace.StartTracing(conf.Tracer, conf.ServiceName); err != nil {
		return err
	}

	// Подключаемся к базе данных
	log.Info("Подключаемся к БД")
	pgsql, err := pgsql.NewClientPgsql(ctx, conf.Pgsql)
	if err != nil {
		return err
	}
	defer pgsql.Close()

	// Запускаем миграции в базе данных
	// TODO: Подумать, как откатывать миграции при ошибках
	log.Info("Запускаем миграции")
	postgreSQLMigrator, err := migrator.NewMigrator(
		migrator.MigratorConfig{
			Migrations:      nil,
			Conn:            pgsql.DB.DB,
			EmbedMigrations: pgsqlMigrations.EmbedMigrationsPgsql,
			Dir:             "pgsql",
			Dialect:         goose.DialectPostgres,
		},
	)
	if err != nil {
		return err
	}
	if err = postgreSQLMigrator.Up(ctx); err != nil {
		return err
	}

	// Регистрируем репозитории
	transactor := transactor.NewTransactor(pgsql)
	accountGroupRepository := accountGroupRepository.NewAccountGroupRepository(pgsql)
	accountRepository := accountRepository.NewAccountRepository(pgsql)
	tagRepository := tagRepository.NewTagRepository(pgsql)
	transactionRepository := transactionRepository.NewTransactionRepository(pgsql)
	settingsRepository := settingsRepository.NewSettingsRepository(pgsql)
	userRepository := userRepository.NewUserRepository(pgsql)
	accountPermissionsRepository := accountPermisssionsRepository.NewAccountPermissionsRepository(pgsql)

	// Регистрируем сервисы
	log.Info("Инициализируем Telegram-бота")
	tgBotService, err := service.NewTgBotService()
	if err != nil {
		return err
	}
	if conf.Telegram.IsEnabled {
		defer tgBotService.Bot.Close()
	}

	log.Info("Инициализируем сервис пушей")
	pushNotificatorService, err := pushNotificatorService.NewPushNotificatorService(conf.Notifications.Enabled, pushNotificatorModel.APNsCredentials{
		TeamID:      conf.Notifications.APNs.TeamID,
		KeyID:       conf.Notifications.APNs.KeyID,
		KeyFilePath: conf.Notifications.APNs.KeyFilePath,
	})
	if err != nil {
		return err
	}

	accountPermissionsService := accountPermisssionsService.NewAccountPermissionsService(accountPermissionsRepository)

	userService := userService.NewUserService(
		userRepository,
		transactor,
		pushNotificatorService,
		[]byte(conf.Auth.GeneralSalt),
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
			CurrencyProviderAPIKey: conf.APIKeys.CurrencyProvider,
		},
	)

	authService := authService.NewAuthService(
		userRepository,
		transactor,
		[]byte(conf.Auth.GeneralSalt),
	)

	log.Info("Запускаем планировщик")
	if err = scheduler.NewScheduler(settingsService).Start(ctx); err != nil {
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

	server, err := server.GetDefaultServer(conf.HTTP, r)
	if err != nil {
		return err
	}

	// Создаем wait группу
	eg, ctx := errgroup.WithContext(ctx)

	// Запускаем HTTP-сервер
	eg.Go(func() error {

		log.Info(fmt.Sprintf("Server is listening: %s", conf.HTTP))

		return server.Serve()
	})

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

func initSingletons(conf config.Config) error {

	// stackTrace.Init(conf.ServiceName, conf.StackTraceEnabled)

	// Конфигурируем decimal, чтобы в JSON не было кавычек
	decimal.MarshalJSONWithoutQuotes = true

	// Инициализируем JWT-менеджер
	accessTokenTTL, err := time.ParseDuration(conf.Auth.AccessTokenTTL)
	if err != nil {
		return errors.InternalServer.Wrap(err)
	}
	refreshTokenTTL, err := time.ParseDuration(conf.Auth.RefreshTokenTTL)
	if err != nil {
		return errors.InternalServer.Wrap(err)
	}
	jwtManager.Init([]byte(conf.Auth.SigningKey), accessTokenTTL, refreshTokenTTL)

	// if err = middleware.(conf.ServiceName); err != nil {
	//	return err
	// }

	return nil
}
