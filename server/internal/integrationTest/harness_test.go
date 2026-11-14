package integrationTest

import (
	"context"
	"sync"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"pkg/passwordManager"
	accountRepository "server/internal/modules/account/repository"
	accountService "server/internal/modules/account/service"
	accountBudgetRepository "server/internal/modules/accountBudget/repository"
	accountBudgetService "server/internal/modules/accountBudget/service"
	accountGroupRepository "server/internal/modules/accountGroup/repository"
	accountGroupService "server/internal/modules/accountGroup/service"
	auditLogRepository "server/internal/modules/auditLog/repository"
	auditLogService "server/internal/modules/auditLog/service"
	authService "server/internal/modules/auth/service"
	pushNotificatorModel "server/internal/modules/pushNotificator/model"
	pushNotificatorService "server/internal/modules/pushNotificator/service"
	settingsRepository "server/internal/modules/settings/repository"
	settingsService "server/internal/modules/settings/service"
	statisticService "server/internal/modules/statistic/service"
	tagRepository "server/internal/modules/tag/repository"
	tagService "server/internal/modules/tag/service"
	tgBotModel "server/internal/modules/tgBot/model"
	transactionRepository "server/internal/modules/transaction/repository"
	transactionService "server/internal/modules/transaction/service"
	"server/internal/modules/transactor"
	userRepository "server/internal/modules/user/repository"
	userService "server/internal/modules/user/service"
	userToAccountGroupRepository "server/internal/modules/userToAccountGroup/repository"
	userToAccountGroupService "server/internal/modules/userToAccountGroup/service"
	"server/internal/utils/necessary"
	"server/internal/utils/syncNotifier"
)

// generalSalt — общая соль паролей тестового сервера
var generalSalt = []byte("salt")

// services — сервисы, связанные так же, как в main.go
type services struct {
	auth         *authService.AuthService                 // Авторизация
	user         *userService.UserService                 // Пользователи и устройства
	account      *accountService.AccountService           // Счета
	accountGroup *accountGroupService.AccountGroupService // Группы счетов
	tag          *tagService.TagService                   // Теги
	settings     *settingsService.SettingsService         // Справочники
	statistic    *statisticService.StatisticService       // Статистика
}

// fixture — окружение одного теста: сервисы и отдельный пользователь с устройством и группой счетов
type fixture struct {
	services
	t              *testing.T                         // Текущий тест
	ctx            context.Context                    // Контекст запросов
	necessary      necessary.NecessaryUserInformation // Пользователь и устройство
	email          string                             // Почта пользователя
	password       string                             // Пароль пользователя
	accountGroupID uuid.UUID                          // Группа счетов пользователя
	iconID         uuid.UUID                          // Иконка для счетов
}

// newServices собирает сервисы поверх тестовой базы
func newServices(t *testing.T) services {
	transactor := transactor.NewTransactor(database)
	accountRepository := accountRepository.NewAccountRepository(database)
	transactionRepository := transactionRepository.NewTransactionRepository(database)
	tagRepository := tagRepository.NewTagRepository(database)
	userRepository := userRepository.NewUserRepository(database)

	pushNotificatorService, err := pushNotificatorService.NewPushNotificatorService(false, pushNotificatorModel.APNsCredentials{}) //nolint:exhaustruct
	if err != nil {
		t.Fatal(err)
	}
	userToAccountGroupService := userToAccountGroupService.NewUserToAccountGroupService(userToAccountGroupRepository.NewUserToAccountGroupRepository(database))
	auditLogService := auditLogService.NewAuditLogService(auditLogRepository.NewAuditLogRepository(database), userToAccountGroupService, syncNotifier.New())
	userService, err := userService.NewUserService(userRepository, transactor, pushNotificatorService, generalSalt, auditLogService)
	if err != nil {
		t.Fatal(err)
	}
	accountGroupService := accountGroupService.NewAccountGroupService(accountGroupRepository.NewAccountGroupRepository(database), transactor, userToAccountGroupService, auditLogService, userService)
	accountService := accountService.NewAccountService(accountRepository, transactor, transactionRepository, userRepository, accountGroupService, userToAccountGroupService, auditLogService, userService)
	accountBudgetService := accountBudgetService.NewAccountBudgetService(accountBudgetRepository.NewAccountBudgetRepository(database), transactor, accountService, accountGroupService, userToAccountGroupService, auditLogService, userService)
	tagService := tagService.NewTagService(tagRepository, transactor, userToAccountGroupService, accountGroupService, auditLogService, userService)
	transactionService := transactionService.NewTransactionService(transactionRepository, accountRepository, transactor, tagRepository, userToAccountGroupService, accountService, tagService, auditLogService, userService)
	settingsService := settingsService.NewSettingsService(settingsRepository.NewSettingsRepository(database), userService, offlineTgBot{}, transactor, auditLogService,
		settingsService.Version{Version: "test", Build: "test"}, settingsService.Credentials{CurrencyProviderAPIKey: ""})
	authService := authService.NewAuthService(userRepository, transactor, generalSalt, 5*time.Minute, auditLogService)

	return services{
		auth:         authService,
		user:         userService,
		account:      accountService,
		accountGroup: accountGroupService,
		tag:          tagService,
		settings:     settingsService,
		statistic: statisticService.NewStatisticService(transactionService, tagService, accountService, accountGroupService,
			userToAccountGroupService, settingsService, userService, accountBudgetService),
	}
}

// offlineTgBot — Telegram-бот, который никуда не пишет: настоящий читает конфиг из окружения
type offlineTgBot struct{}

// SendMessage ничего не отправляет
func (offlineTgBot) SendMessage(context.Context, tgBotModel.SendMessageReq) error { return nil }

// ratesOnce — курсы валют выставляются один раз на всю базу
var ratesOnce sync.Once

// newFixture заводит нового пользователя с паролем, устройством и группой счетов в рублях
func newFixture(t *testing.T) *fixture {
	ctx := context.Background()
	f := &fixture{
		services:       newServices(t),
		t:              t,
		ctx:            ctx,
		necessary:      necessary.NecessaryUserInformation{UserID: uuid.New(), DeviceID: "test-" + uuid.NewString()[:20]},
		password:       "secret-password",
		accountGroupID: uuid.New(),
	}
	f.email = f.necessary.UserID.String()[:18] + "@t.local"

	// Детерминированные курсы в формате базы — единиц валюты за 1 USD: 100 RUB, 0.00001 BTC
	ratesOnce.Do(func() {
		f.exec(sq.Update("coin.currencies").Set("rate", sq.Expr("CASE slug WHEN 'USD' THEN 1 WHEN 'RUB' THEN 100 WHEN 'BTC' THEN 0.00001 ELSE rate END")))

		// История: те же курсы действуют с 2000 года, а EUR меняется — для теста переоценки:
		// до октября 2026 года 1 EUR = 1 USD, с октября 1 EUR = 2 USD
		f.exec(sq.Delete("coin.currency_rates").Where(sq.Eq{"currency": []string{"USD", "RUB", "BTC", "EUR"}}))
		points := sq.Insert("coin.currency_rates").Columns("currency", "valid_at", "rate", "source")
		for currency, rate := range map[string]string{"USD": "1", "RUB": "100", "BTC": "0.00001", "EUR": "1"} {
			points = points.Values(currency, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), rate, "provider")
		}
		f.exec(points.Values("EUR", time.Date(2026, 9, 30, 23, 59, 59, 0, time.UTC).Add(time.Second), "0.5", "provider"))
	})

	// Пользователь с настоящим хэшем пароля и устройство
	userSalt, err := passwordManager.GenerateRandomSalt()
	f.ok(err)
	hash, err := passwordManager.CreateNewPassword([]byte(f.password), generalSalt, userSalt)
	f.ok(err)
	f.exec(sq.Insert("coin.users").SetMap(map[string]any{
		"id": f.necessary.UserID, "name": "Test", "email": f.email, "password_hash": hash, "password_salt": userSalt,
		"time_create": time.Now(), "default_currency": "USD",
	}))
	f.exec(sq.Insert("coin.devices").SetMap(map[string]any{
		"user_id": f.necessary.UserID, "device_id": f.necessary.DeviceID, "refresh_token": "token",
		"application_bundle_id": "test", "application_version": "1", "application_build": "1",
		"device_ip_address": "127.0.0.1", "device_user_agent": "test", "device_os_name": "iOS",
		"device_os_version": "1", "device_name": "test", "device_model_name": "test",
	}))

	// Группа счетов в рублях
	f.exec(sq.Insert("coin.account_groups").SetMap(map[string]any{
		"id": f.accountGroupID, "name": "Личные", "currency_signatura": "RUB", "serial_number": 1, "visible": true,
		"datetime_create": time.Now(), "created_by_user_id": f.necessary.UserID,
	}))
	f.exec(sq.Insert("coin.users_to_account_groups").SetMap(map[string]any{"user_id": f.necessary.UserID, "account_group_id": f.accountGroupID}))

	// Иконка для счетов
	f.ok(database.Get(ctx, &f.iconID, sq.Select("id").From("coin.icons").Limit(1)))
	return f
}

// accountOptions — необязательные поля счета
type accountOptions struct {
	parentID *uuid.UUID // Родительский счет
	isParent bool       // Счет-родитель
}

// account создает счет в группе пользователя
func (f *fixture) newAccount(name, accountType, currency string, options accountOptions) uuid.UUID {
	f.t.Helper()
	id := uuid.New()
	f.exec(sq.Insert("coin.accounts").SetMap(map[string]any{
		"id": id, "name": name, "account_type": accountType, "currency": currency, "rank": name,
		"datetime_create": time.Now(), "account_group_id": f.accountGroupID, "created_by_user_id": f.necessary.UserID,
		"icon_id": f.iconID, "parent_account_id": options.parentID, "is_parent": options.isParent,
	}))
	return id
}

// tag создает тег в группе пользователя
func (f *fixture) newTag(name string) uuid.UUID {
	f.t.Helper()
	id := uuid.New()
	f.exec(sq.Insert("coin.tags").SetMap(map[string]any{
		"id": id, "name": name, "datetime_create": time.Now(), "account_group_id": f.accountGroupID, "created_by_user_id": f.necessary.UserID,
	}))
	return id
}

// transaction создает транзакцию с тегами
func (f *fixture) newTransaction(transactionType string, from, to uuid.UUID, amountFrom, amountTo string, date string, note string, tags ...uuid.UUID) uuid.UUID {
	f.t.Helper()
	id := uuid.New()
	f.exec(sq.Insert("coin.transactions").SetMap(map[string]any{
		"id": id, "date_transaction": date, "transaction_type": transactionType,
		"amount_from": decimal.RequireFromString(amountFrom), "amount_to": decimal.RequireFromString(amountTo), "note": note,
		"datetime_create": time.Now(), "account_from_id": from, "account_to_id": to,
		"created_by_user_id": f.necessary.UserID, "account_group_id": f.accountGroupID,
	}))
	for _, tagID := range tags {
		f.exec(sq.Insert("coin.tags_to_transaction").SetMap(map[string]any{"tag_id": tagID, "transaction_id": id}))
	}
	return id
}

// budget создает версию бюджета расходного счета
func (f *fixture) newBudget(accountID uuid.UUID, amount string, effectiveFrom string) {
	f.t.Helper()
	f.exec(sq.Insert("coin.account_budgets").SetMap(map[string]any{
		"id": uuid.New(), "account_id": accountID, "amount": decimal.RequireFromString(amount), "effective_from": effectiveFrom,
		"created_by_user_id": f.necessary.UserID, "account_group_id": f.accountGroupID,
	}))
}

// exec выполняет служебный запрос подготовки данных
func (f *fixture) exec(query sq.Sqlizer) {
	f.t.Helper()
	if err := database.Exec(f.ctx, query); err != nil {
		f.t.Fatal(err)
	}
}

// ok падает на неожиданной ошибке
func (f *fixture) ok(err error) {
	f.t.Helper()
	if err != nil {
		f.t.Fatalf("неожиданная ошибка: %v", err)
	}
}

// count возвращает число строк по запросу COUNT(*)
func (f *fixture) count(query sq.SelectBuilder) int {
	f.t.Helper()
	var result int
	if err := database.Get(f.ctx, &result, query); err != nil {
		f.t.Fatal(err)
	}
	return result
}
