// import-rates загружает историю курсов валют из currencyapi в coin.currency_rates.
//
// Использование:
//
//	go run ./cmd/import-rates -from 2019-10-31 -to 2026-11-13
//
// Один день — один запрос /v3/historical со всеми валютами. Дни, за которые курсы провайдера уже
// есть, пропускаются, поэтому команду можно перезапускать: например, когда закончилась квота.
// Параметры из env: PGSQL_HOST, PGSQL_DATABASE, PGSQL_USER, PGSQL_PASSWORD и
// CURRENCY_API_KEY — ключ с достаточной квотой (около 2 600 запросов на всю историю).
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pkg/database/pgsql"
	settingsNetwork "server/internal/modules/settings/network"
	settingsRepository "server/internal/modules/settings/repository"
)

// dateLayout — формат дат флагов
const dateLayout = "2006-01-02"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(1)
	}
}

// run загружает недостающие дни истории
func run() error {
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Format(dateLayout)
	from := flag.String("from", "2019-10-31", "первый день, yyyy-MM-dd")
	to := flag.String("to", yesterday, "последний день включительно, yyyy-MM-dd")
	delay := flag.Duration("delay", 500*time.Millisecond, "пауза между запросами к провайдеру")
	flag.Parse()

	fromDate, err := time.Parse(dateLayout, *from)
	if err != nil {
		return fmt.Errorf("-from: %w", err)
	}
	toDate, err := time.Parse(dateLayout, *to)
	if err != nil {
		return fmt.Errorf("-to: %w", err)
	}
	apiKey := os.Getenv("CURRENCY_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("не задан CURRENCY_API_KEY")
	}

	// Останавливаемся по Ctrl+C или при остановке Job: уже загруженные дни сохранены
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Подключаемся к базе
	database, err := pgsql.NewClientPgsql(ctx, pgsql.PgsqlConfigEnv{
		Host:     os.Getenv("PGSQL_HOST"),
		User:     os.Getenv("PGSQL_USER"),
		Password: os.Getenv("PGSQL_PASSWORD"),
		Database: os.Getenv("PGSQL_DATABASE"),
	})
	if err != nil {
		return err
	}
	defer database.Close()
	repository := settingsRepository.NewSettingsRepository(database)

	// Дни, которые уже загружены
	loaded, err := repository.GetProviderRateDays(ctx, fromDate, toDate.AddDate(0, 0, 1))
	if err != nil {
		return err
	}

	// Загружаем недостающие дни по порядку
	imported, skipped := 0, 0
	for day := fromDate; !day.After(toDate); day = day.AddDate(0, 0, 1) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if loaded[day.Format(dateLayout)] {
			skipped++
			continue
		}

		rates, validAt, err := settingsNetwork.GetHistoricalCurrencyRates(ctx, apiKey, day)
		if err != nil {
			return fmt.Errorf("%s: %w (загружено %d, пропущено %d — перезапуск продолжит с этого дня)", day.Format(dateLayout), err, imported, skipped)
		}

		// Провайдер считает исторический курс на конец дня; если момента нет — ставим конец дня сами
		if validAt.IsZero() {
			validAt = day.AddDate(0, 0, 1).Add(-time.Second)
		}
		if err = repository.InsertCurrencyRates(ctx, validAt, rates); err != nil {
			return fmt.Errorf("%s: %w", day.Format(dateLayout), err)
		}
		imported++
		if imported%50 == 0 {
			fmt.Printf("%s: загружено %d дней\n", day.Format(dateLayout), imported)
		}
		time.Sleep(*delay)
	}

	fmt.Printf("готово: загружено %d дней, пропущено %d уже загруженных\n", imported, skipped)
	return nil
}
