package integrationTest

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/google/uuid"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	accountEndpointMCP "server/internal/modules/account/endpoint/mcp"
	accountGroupEndpointMCP "server/internal/modules/accountGroup/endpoint/mcp"
	settingsEndpointMCP "server/internal/modules/settings/endpoint/mcp"
	statisticEndpointMCP "server/internal/modules/statistic/endpoint/mcp"
	tagEndpointMCP "server/internal/modules/tag/endpoint/mcp"
	"server/internal/utils/contextKeys"
	mcpTransport "server/internal/utils/mcp"
)

// newMCPServer собирает MCP-сервер со всеми инструментами, как в main.go
func (f *fixture) newMCPServer() *sdk.Server {
	server := mcpTransport.NewServer("test")
	statisticEndpointMCP.Register(server, f.statistic)
	accountEndpointMCP.Register(server, f.account)
	tagEndpointMCP.Register(server, f.tag)
	accountGroupEndpointMCP.Register(server, f.accountGroup)
	settingsEndpointMCP.Register(server, f.settings)
	return server
}

// mcpSession подключает клиента к MCP-серверу через in-memory транспорт. Пользователь и устройство —
// в контексте сервера, как после проверки Bearer-токена
func (f *fixture) mcpSession() *sdk.ClientSession {
	f.t.Helper()
	ctx := contextKeys.SetUserID(f.ctx, f.necessary.UserID)
	ctx = contextKeys.SetDeviceID(ctx, f.necessary.DeviceID)

	serverTransport, clientTransport := sdk.NewInMemoryTransports()
	if _, err := f.newMCPServer().Connect(ctx, serverTransport, nil); err != nil {
		f.t.Fatal(err)
	}
	session, err := sdk.NewClient(&sdk.Implementation{Name: "test", Version: "0"}, nil).Connect(ctx, clientTransport, nil) //nolint:exhaustruct
	if err != nil {
		f.t.Fatal(err)
	}
	f.t.Cleanup(func() { _ = session.Close() })
	return session
}

// call вызывает инструмент и разбирает структурированный ответ; ошибка инструмента — в тексте
func (f *fixture) call(session *sdk.ClientSession, name string, args map[string]any, out any) error {
	f.t.Helper()
	result, err := session.CallTool(f.ctx, &sdk.CallToolParams{Name: name, Arguments: args}) //nolint:exhaustruct
	if err != nil {
		f.t.Fatal(err)
	}
	if result.IsError {
		text, _ := result.Content[0].(*sdk.TextContent)
		return &toolError{text: text.Text}
	}
	if out != nil {
		raw, _ := json.Marshal(result.StructuredContent)
		if err := json.Unmarshal(raw, out); err != nil {
			f.t.Fatal(err)
		}
	}
	return nil
}

// toolError — ошибка, которую инструмент вернул агенту текстом
type toolError struct{ text string }

func (e *toolError) Error() string { return e.text }

// statisticOut — ответ get_statistic с суммами строками, чтобы сравнивать без float
type statisticOut struct {
	Currency string   `json:"currency"`
	Buckets  []string `json:"buckets"`
	Series   []struct {
		Key            string        `json:"key"`
		Name           string        `json:"name"`
		Values         []json.Number `json:"values"`
		Sum            json.Number   `json:"sum"`
		Max            json.Number   `json:"max"`
		AverageNonZero json.Number   `json:"averageNonZero"`
		MonthlyBudget  *json.Number  `json:"monthlyBudget"`
	} `json:"series"`
}

// values возвращает значения ряда по названию строками
func (s statisticOut) values(t *testing.T, name string) string {
	t.Helper()
	for _, series := range s.Series {
		if series.Name == name {
			parts := make([]string, 0, len(series.Values))
			for _, value := range series.Values {
				parts = append(parts, value.String())
			}
			return strings.Join(parts, " ")
		}
	}
	t.Fatalf("нет ряда %q в %+v", name, s.Series)
	return ""
}

// world — счета и теги сценария
type world struct {
	card, wallet, salary, food, groceries, cafe, taxi, balancing uuid.UUID
	vegetables, coffee                                           uuid.UUID
}

// seed заводит сценарий: зарплата, расходы по иерархии категорий и тегам, валютный перевод,
// корректировка остатка и доход до начала периода
func (f *fixture) seed() world {
	w := world{}
	w.card = f.newAccount("Карта", "regular", "RUB", accountOptions{})
	w.wallet = f.newAccount("Кошелек USD", "regular", "USD", accountOptions{})
	w.salary = f.newAccount("Зарплата", "earnings", "RUB", accountOptions{})
	w.food = f.newAccount("Еда", "expense", "RUB", accountOptions{isParent: true})
	w.groceries = f.newAccount("Продукты", "expense", "RUB", accountOptions{parentID: &w.food})
	w.cafe = f.newAccount("Кафе", "expense", "RUB", accountOptions{parentID: &w.food})
	w.taxi = f.newAccount("Такси", "expense", "RUB", accountOptions{})
	w.balancing = f.newAccount("Корректировка", "balancing", "RUB", accountOptions{})
	w.vegetables = f.newTag("Овощи")
	w.coffee = f.newTag("Кофе")

	f.newTransaction("income", w.salary, w.card, "1000", "1000", "2026-08-20", "аванс")
	f.newTransaction("income", w.salary, w.card, "100000", "100000", "2026-09-01", "зарплата")
	f.newTransaction("consumption", w.card, w.groceries, "3000", "3000", "2026-09-05", "Рынок", w.vegetables)
	f.newTransaction("consumption", w.card, w.cafe, "2000", "2000", "2026-10-10", "капучино", w.coffee)
	f.newTransaction("consumption", w.wallet, w.taxi, "10", "1000", "2026-10-12", "")
	f.newTransaction("transfer", w.card, w.wallet, "5000", "50", "2026-10-15", "обмен")
	f.newTransaction("balancing", w.balancing, w.card, "500", "500", "2026-11-01", "")
	f.newBudget(w.groceries, "5000", "2026-09-01")
	return w
}

// statistic вызывает get_statistic за сентябрь–ноябрь 2026 по месяцам
func (f *fixture) statisticCall(session *sdk.ClientSession, args map[string]any) statisticOut {
	f.t.Helper()
	args["period"] = "month"
	filters, _ := args["filters"].(map[string]any)
	if filters == nil {
		filters = map[string]any{}
	}
	filters["dateFrom"], filters["dateTo"] = "2026-09-01", "2026-11-30"
	args["filters"] = filters
	var out statisticOut
	f.ok(f.call(session, "get_statistic", args, &out))
	return out
}

func TestMCP_ToolsListed(t *testing.T) {
	f := newFixture(t)
	tools, err := f.mcpSession().ListTools(f.ctx, nil)
	f.ok(err)
	names := make([]string, 0, len(tools.Tools))
	for _, tool := range tools.Tools {
		names = append(names, tool.Name)
		if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
			t.Fatalf("инструмент %s не помечен только на чтение", tool.Name)
		}
	}
	slices.Sort(names)
	expected := []string{"get_account_groups", "get_accounts", "get_currencies", "get_statistic", "get_tags", "get_transactions"}
	if !slices.Equal(names, expected) {
		t.Fatalf("инструменты: %v", names)
	}
}

func TestMCP_Statistic(t *testing.T) {
	f := newFixture(t)
	w := f.seed()
	session := f.mcpSession()

	// Доходы и расходы в валюте единственной группы — рублях
	out := f.statisticCall(session, map[string]any{"chartType": "income_expense"})
	if out.Currency != "RUB" || strings.Join(out.Buckets, " ") != "2026-09-01 2026-10-01 2026-11-01" {
		t.Fatalf("валюта и бакеты: %s %v", out.Currency, out.Buckets)
	}
	if got := out.values(t, "Доходы"); got != "100000 0 0" {
		t.Fatalf("доходы: %s", got)
	}
	if got := out.values(t, "Расходы"); got != "3000 3000 0" {
		t.Fatalf("расходы: %s", got)
	}

	// Расходы по категориям: дети свернуты в родителя с isParent, у него бюджет детей
	out = f.statisticCall(session, map[string]any{"chartType": "expense"})
	if got := out.values(t, "Еда"); got != "3000 2000 0" {
		t.Fatalf("еда: %s", got)
	}
	if got := out.values(t, "Такси"); got != "0 1000 0" {
		t.Fatalf("такси: %s", got)
	}
	if out.Series[0].Name != "Еда" || out.Series[0].Sum.String() != "5000" || out.Series[0].AverageNonZero.String() != "2500" ||
		out.Series[0].MonthlyBudget == nil || out.Series[0].MonthlyBudget.String() != "5000" {
		t.Fatalf("сводка еды: %+v", out.Series[0])
	}

	// Без сворачивания — дочерние категории отдельно
	out = f.statisticCall(session, map[string]any{"chartType": "expense", "aggregateIntoParents": false})
	if got := out.values(t, "Продукты"); got != "3000 0 0" {
		t.Fatalf("продукты: %s", got)
	}
	if got := out.values(t, "Кафе"); got != "0 2000 0" {
		t.Fatalf("кафе: %s", got)
	}

	// По тегам: без тега — отдельный ряд
	out = f.statisticCall(session, map[string]any{"chartType": "expense", "groupBy": "tag"})
	if got := out.values(t, "Овощи") + " | " + out.values(t, "Кофе") + " | " + out.values(t, "Без тега"); got != "3000 0 0 | 0 2000 0 | 0 1000 0" {
		t.Fatalf("теги: %s", got)
	}

	// Дельта и пересчет в выбранную нейросетью валюту
	out = f.statisticCall(session, map[string]any{"chartType": "delta"})
	if got := out.values(t, "Доходы минус расходы"); got != "97000 -3000 0" {
		t.Fatalf("дельта: %s", got)
	}
	out = f.statisticCall(session, map[string]any{"chartType": "income_expense", "currency": "USD"})
	if got := out.values(t, "Расходы"); out.Currency != "USD" || got != "30 30 0" {
		t.Fatalf("расходы в USD: %s %s", out.Currency, got)
	}

	// В наличии от первой транзакции: аванс до периода входит во входящий остаток, корректировка — в ноябрь
	out = f.statisticCall(session, map[string]any{"chartType": "balance_detailed"})
	if got := out.values(t, "Карта"); got != "98000 91000 91500" {
		t.Fatalf("карта: %s", got)
	}
	if got := out.values(t, "Кошелек USD"); got != "0 4000 4000" {
		t.Fatalf("кошелек: %s", got)
	}
	out = f.statisticCall(session, map[string]any{"chartType": "balance_total"})
	if got := out.values(t, "В наличии"); got != "98000 95000 95500" {
		t.Fatalf("всего: %s", got)
	}

	// Фильтр счета выбирает счета наличности
	out = f.statisticCall(session, map[string]any{"chartType": "balance_total", "filters": map[string]any{"accountIDs": []string{w.wallet.String()}}})
	if got := out.values(t, "В наличии"); got != "0 4000 4000" {
		t.Fatalf("всего по кошельку: %s", got)
	}

	// Ошибки параметров возвращаются агенту текстом
	if err := f.call(session, "get_statistic", map[string]any{"chartType": "pie", "period": "month"}, nil); err == nil {
		t.Fatal("неизвестный chartType принят")
	}
}

// transactionsOut — ответ get_transactions
type transactionsOut struct {
	Currency         string      `json:"currency"`
	TotalCount       int         `json:"totalCount"`
	TotalConsumption json.Number `json:"totalConsumption"`
	TotalIncome      json.Number `json:"totalIncome"`
	NextOffset       *int        `json:"nextOffset"`
	Days             []struct {
		Date             string      `json:"date"`
		ConsumptionTotal json.Number `json:"consumptionTotal"`
		Transactions     []struct {
			Note          string   `json:"note"`
			AccountToName string   `json:"accountToName"`
			TagIDs        []string `json:"tagIDs"`
		} `json:"transactions"`
	} `json:"days"`
}

func TestMCP_Transactions(t *testing.T) {
	f := newFixture(t)
	w := f.seed()
	session := f.mcpSession()

	// Октябрь: три транзакции по дням от новых к старым, расходы за день и по выборке
	var out transactionsOut
	f.ok(f.call(session, "get_transactions", map[string]any{"filters": map[string]any{"dateFrom": "2026-10-01", "dateTo": "2026-10-31"}}, &out))
	days := make([]string, 0, len(out.Days))
	for _, day := range out.Days {
		days = append(days, day.Date+"="+day.ConsumptionTotal.String())
	}
	if out.TotalCount != 3 || out.TotalConsumption.String() != "3000" || strings.Join(days, " ") != "2026-10-15=0 2026-10-12=1000 2026-10-10=2000" {
		t.Fatalf("октябрь: %d %s %v", out.TotalCount, out.TotalConsumption, days)
	}

	// Пагинация: итоги по всей выборке, страница — limit
	f.ok(f.call(session, "get_transactions", map[string]any{"limit": 2}, &out))
	if out.TotalCount != 7 || out.NextOffset == nil || *out.NextOffset != 2 || out.TotalIncome.String() != "101000" {
		t.Fatalf("страница: %d %v %s", out.TotalCount, out.NextOffset, out.TotalIncome)
	}
	out = transactionsOut{}
	f.ok(f.call(session, "get_transactions", map[string]any{"limit": 500, "offset": 6}, &out))
	if out.NextOffset != nil || len(out.Days) != 1 || out.Days[0].Transactions[0].Note != "аванс" {
		t.Fatalf("последняя страница: %+v", out)
	}

	// Фильтры: сумма в валюте счета, заметка, родительская категория с детьми, исключенный тег
	cases := []struct {
		name    string
		filters map[string]any
		count   int
	}{
		{"сумма от 2500", map[string]any{"amountMin": 2500}, 3},
		{"сумма до 10", map[string]any{"amountMax": 10}, 1},
		{"заметка", map[string]any{"searchText": "КАПУ"}, 1},
		{"категория с детьми", map[string]any{"accountIDs": []string{w.food.String()}}, 2},
		{"исключенный счет", map[string]any{"excludedAccountIDs": []string{w.card.String()}}, 1},
		{"тег", map[string]any{"tagIDs": []string{w.coffee.String()}}, 1},
		{"исключенный тег", map[string]any{"excludedTagIDs": []string{w.coffee.String(), w.vegetables.String()}}, 5},
		{"тип", map[string]any{"types": []string{"consumption"}}, 3},
		{"валюта", map[string]any{"currencies": []string{"USD"}}, 2},
	}
	for _, c := range cases {
		f.ok(f.call(session, "get_transactions", map[string]any{"filters": c.filters}, &out))
		if out.TotalCount != c.count {
			t.Fatalf("%s: %d вместо %d", c.name, out.TotalCount, c.count)
		}
	}
}

func TestMCP_AccessIsolated(t *testing.T) {
	owner := newFixture(t)
	owner.seed()
	stranger := newFixture(t)
	session := stranger.mcpSession()

	// Чужие транзакции не видны
	var out transactionsOut
	stranger.ok(stranger.call(session, "get_transactions", map[string]any{}, &out))
	if out.TotalCount != 0 {
		t.Fatalf("видны чужие транзакции: %d", out.TotalCount)
	}

	// Чужая группа счетов в фильтре — отказ
	err := stranger.call(session, "get_statistic", map[string]any{"chartType": "expense", "period": "month",
		"filters": map[string]any{"accountGroupIDs": []string{owner.accountGroupID.String()}}}, nil)
	if err == nil {
		t.Fatal("чужая группа счетов принята")
	}

	// Справочники — только свои
	var groups struct {
		AccountGroups []struct {
			ID string `json:"id"`
		} `json:"accountGroups"`
	}
	stranger.ok(stranger.call(session, "get_account_groups", map[string]any{}, &groups))
	if len(groups.AccountGroups) != 1 || groups.AccountGroups[0].ID != stranger.accountGroupID.String() {
		t.Fatalf("группы: %+v", groups)
	}
}

func TestMCP_Directories(t *testing.T) {
	f := newFixture(t)
	w := f.seed()
	session := f.mcpSession()

	// Счета: остаток денег текущий, у категорий — сумма за период
	var accounts struct {
		Accounts []struct {
			ID              string      `json:"id"`
			Remainder       json.Number `json:"remainder"`
			ParentAccountID *string     `json:"parentAccountID"`
		} `json:"accounts"`
	}
	f.ok(f.call(session, "get_accounts", map[string]any{"dateFrom": "2026-10-01", "dateTo": "2026-10-31"}, &accounts))
	remainders := map[string]string{}
	for _, account := range accounts.Accounts {
		remainders[account.ID] = account.Remainder.String()
	}
	if remainders[w.card.String()] != "91500" || remainders[w.cafe.String()] != "2000" || len(accounts.Accounts) != 8 {
		t.Fatalf("счета: %v", remainders)
	}

	// Теги и валюты
	var tags struct {
		Tags []struct {
			Name string `json:"name"`
		} `json:"tags"`
	}
	f.ok(f.call(session, "get_tags", map[string]any{}, &tags))
	if len(tags.Tags) != 2 {
		t.Fatalf("теги: %+v", tags)
	}
	var currencies struct {
		Currencies []struct {
			Code string      `json:"code"`
			Rate json.Number `json:"rate"`
		} `json:"currencies"`
	}
	f.ok(f.call(session, "get_currencies", map[string]any{}, &currencies))
	if !slices.ContainsFunc(currencies.Currencies, func(c struct {
		Code string      `json:"code"`
		Rate json.Number `json:"rate"`
	}) bool {
		return c.Code == "RUB" && c.Rate.String() == "100"
	}) {
		t.Fatal("нет RUB с курсом 100")
	}
}
