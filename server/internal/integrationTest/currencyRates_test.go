package integrationTest

import (
	settingsModel "server/internal/modules/settings/model"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/shopspring/decimal"

	settingsRepository "server/internal/modules/settings/repository"
)

func TestCurrencyRates_ImportIsIdempotent(t *testing.T) {
	f := newFixture(t)
	repository := settingsRepository.NewSettingsRepository(database)

	// Миграция переносит текущие курсы первой точкой истории
	if n := f.count(sq.Select("count(*)").From("coin.currency_rates").Where(sq.Eq{"currency": "RUB"})); n < 1 {
		t.Fatalf("текущие курсы не перенесены в историю: %d", n)
	}

	// Курсы провайдера на конец дня; повтор того же момента перезаписывает курс, а не дублирует
	day := time.Date(2020, 3, 15, 0, 0, 0, 0, time.UTC)
	validAt := day.Add(24*time.Hour - time.Second)
	f.ok(repository.InsertCurrencyRates(f.ctx, validAt, map[string]decimal.Decimal{"TST": decimal.RequireFromString("74.5"), "ZZZ": decimal.Zero}))
	f.ok(repository.InsertCurrencyRates(f.ctx, validAt, map[string]decimal.Decimal{"TST": decimal.RequireFromString("75")}))
	var rates []string
	f.ok(database.Select(f.ctx, &rates, sq.Select("rate::text").From("coin.currency_rates").Where(sq.Eq{"currency": "TST", "valid_at": validAt})))
	if len(rates) != 1 || rates[0] != "75" {
		t.Fatalf("точки дня: %v", rates)
	}

	// Нулевой курс не сохраняется
	if n := f.count(sq.Select("count(*)").From("coin.currency_rates").Where(sq.Eq{"currency": "ZZZ"})); n != 0 {
		t.Fatalf("сохранен нулевой курс")
	}

	// Загруженный день виден при повторном запуске, соседний — нет
	days, err := repository.GetProviderRateDays(f.ctx, day.AddDate(0, 0, -1), day.AddDate(0, 0, 2))
	f.ok(err)
	if !days["2020-03-15"] || days["2020-03-14"] || days["2020-03-16"] {
		t.Fatalf("загруженные дни: %v", days)
	}
}

func TestCurrencyRates_PageDoesNotSplitMoment(t *testing.T) {
	f := newFixture(t)
	repository := settingsRepository.NewSettingsRepository(database)

	// Три момента по две валюты
	base := time.Date(2019, 1, 1, 23, 59, 59, 0, time.UTC)
	for day := 0; day < 3; day++ {
		f.ok(repository.InsertCurrencyRates(f.ctx, base.AddDate(0, 0, day), map[string]decimal.Decimal{"TSA": decimal.NewFromInt(int64(day + 1)), "TSB": decimal.NewFromInt(int64(day + 10))}))
	}
	req := settingsModel.GetCurrencyRatesReq{Currencies: []string{"TSA", "TSB"}, Limit: 3} //nolint:exhaustruct

	// Страница из трех точек не режет второй момент: в ней только первый, и есть продолжение
	page, err := f.settings.GetCurrencyRatesPage(f.ctx, req)
	f.ok(err)
	if len(page.Rates) != 2 || !page.HasMore {
		t.Fatalf("первая страница: %d %v", len(page.Rates), page.HasMore)
	}

	// Продолжение с последнего момента доходит до конца
	req.Since = &page.Rates[len(page.Rates)-1].ValidAt
	req.Limit = 10
	page, err = f.settings.GetCurrencyRatesPage(f.ctx, req)
	f.ok(err)
	if len(page.Rates) != 4 || page.HasMore || page.Rates[0].Rate.String() != "2" {
		t.Fatalf("продолжение: %+v", page)
	}
}
