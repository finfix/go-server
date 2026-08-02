package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"pkg/gracefulShutdown"
	fiber2 "pkg/http/fiber"
	"pkg/log/model"
	"pkg/middleware"
	"pkg/stackTrace"
	"server/internal/metrics"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/finfix/go-server-grpc/proto"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/pressly/goose/v3"
	"github.com/shopspring/decimal"

	"pkg/database/pgsql"
	grpcServer "pkg/grpc/server"
	"server/internal/utils/grpc/interceptor"

	"pkg/jwtManager"
	"pkg/log"
	"pkg/migrator"
	"pkg/panicRecover"
	"server/internal/config"
	accountEndpointGRPC "server/internal/modules/account/endpoint/grpc"
	accountRepository "server/internal/modules/account/repository"
	accountService "server/internal/modules/account/service"
	accountGroupEndpointGRPC "server/internal/modules/accountGroup/endpoint/grpc"
	accountGroupRepository "server/internal/modules/accountGroup/repository"
	accountGroupService "server/internal/modules/accountGroup/service"
	authEndpointGRPC "server/internal/modules/auth/endpoint/grpc"
	authService "server/internal/modules/auth/service"
	pushNotificatorModel "server/internal/modules/pushNotificator/model"
	pushNotificatorService "server/internal/modules/pushNotificator/service"
	"server/internal/modules/scheduler"
	settingsEndpointGRPC "server/internal/modules/settings/endpoint/grpc"
	settingsRepository "server/internal/modules/settings/repository"
	settingsService "server/internal/modules/settings/service"
	tagEndpointGRPC "server/internal/modules/tag/endpoint/grpc"
	tagRepository "server/internal/modules/tag/repository"
	tagService "server/internal/modules/tag/service"
	"server/internal/modules/tgBot/service"
	transactionEndpointGRPC "server/internal/modules/transaction/endpoint/grpc"
	transactionRepository "server/internal/modules/transaction/repository"
	transactionService "server/internal/modules/transaction/service"
	"server/internal/modules/transactor"
	userEndpointGRPC "server/internal/modules/user/endpoint/grpc"
	userRepository "server/internal/modules/user/repository"
	userService "server/internal/modules/user/service"
	"server/internal/utils/errors"
	pgsqlMigrations "server/migrations/pgsql"
)

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

	userService, err := userService.NewUserService(
		userRepository,
		transactor,
		pushNotificatorService,
		[]byte(conf.Auth.GeneralSalt),
	)
	if err != nil {
		return err
	}

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

	// Настраиваем и получаем GRPC-сервер
	grpcServer := grpcServer.NewGRPCServer(&grpcServer.ServerOptions{
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			interceptor.NewErrorHandlerInterceptor().Unary(),
			interceptor.NewTimeInterceptor().Unary(),
			interceptor.NewAuthInterceptor(config.GetOpenForEveryoneMethods()).Unary(),
			interceptor.NewLoggerInterceptor().Unary(),
		},
		StreamInterceptors: []grpc.StreamServerInterceptor{
			interceptor.NewErrorHandlerInterceptor().Stream(),
			interceptor.NewTimeInterceptor().Stream(),
			interceptor.NewAuthInterceptor(config.GetOpenForEveryoneMethods()).Stream(),
			interceptor.NewLoggerInterceptor().Stream(),
		},
	})
	proto.RegisterAccountEndpointServer(grpcServer, accountEndpointGRPC.NewAccountEndpoint(accountService))
	proto.RegisterAccountGroupEndpointServer(grpcServer, accountGroupEndpointGRPC.NewAccountGroupEndpoint(accountGroupService))
	proto.RegisterTagEndpointServer(grpcServer, tagEndpointGRPC.NewTagEndpoint(tagService))
	proto.RegisterTransactionEndpointServer(grpcServer, transactionEndpointGRPC.NewTransactionEndpoint(transactionService))
	proto.RegisterUserEndpointServer(grpcServer, userEndpointGRPC.NewUserEndpoint(userService))
	proto.RegisterSettingsEndpointServer(grpcServer, settingsEndpointGRPC.NewSettingsEndpoint(settingsService))
	proto.RegisterAuthEndpointServer(grpcServer, authEndpointGRPC.NewAuthEndpoint(authService))

	// Создаем слушателя порта для gRPC-сервера
	grpcLn, err := net.Listen("tcp", conf.Listen.GRPC)
	if err != nil {
		return errors.InternalServer.Wrap(err)
	}
	defer grpcLn.Close()

	// Настраиваем и получаем HTTP-сервер
	httpServer, err := fiber2.GetDefaultServer(conf.ServiceName, nil)
	if err != nil {
		return err
	}

	// Настраиваем и получаем служебный HTTP-сервер
	serviceHTTPServer, err := fiber2.GetDefaultServer(conf.ServiceName, nil)
	if err != nil {
		return err
	}

	// Настраиваем служебный HTTP-сервер для pprof
	serviceHTTPServer.Use(pprof.New())

	// Инициализируем эндпоинты
	httpServer.Get("/version", middleware.NewVersionHandler(version, build, buildDate, hostname)) // GET /version

	// Создаем wait группу
	eg, ctx := errgroup.WithContext(ctx)

	// Создаем горутину на запуск HTTP-сервера
	eg.Go(func() error {
		log.Info(fmt.Sprintf("HTTP server is listening :%v", conf.Listen.HTTP))
		if err := httpServer.Listen(conf.Listen.HTTP); err != nil {
			return errors.InternalServer.Wrap(err)
		}
		return nil
	})

	// Создаем горутину на запуск служебного HTTP-сервера
	eg.Go(func() error {
		log.Info(fmt.Sprintf("service HTTP server is listening :%v", conf.Listen.ServiceHTTP))
		if err := serviceHTTPServer.Listen(conf.Listen.ServiceHTTP); err != nil {
			return errors.InternalServer.Wrap(err)
		}
		return nil
	})

	// Создаем горутину на запуск gRPC-сервера
	eg.Go(func() error {
		log.Info(fmt.Sprintf("gRPC server is listening: %v", conf.Listen.GRPC))
		if err := grpcServer.Serve(grpcLn); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			return errors.InternalServer.Wrap(err)
		}
		return nil
	})

	// Добавляем в группу горутину для graceful shutdown
	gracefulShutdown.AddGracefulShutdownErrGroup(eg, ctx, []*fiber.App{httpServer, serviceHTTPServer}, []*grpc.Server{grpcServer})

	// Ждем завершения контекста или ошибок в горутинах
	return eg.Wait()
}

func initSingletons(conf config.Config) error {

	stackTrace.Init(conf.ServiceName, true, 1)

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

	return metrics.Init(conf.ServiceName)
}
