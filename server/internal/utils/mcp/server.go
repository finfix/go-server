package mcp

import (
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// instructions — что агенту нужно знать об API. Только семантика полей, без правил работы с
// данными: они принадлежат пользователю и приходят из его инструкций
const instructions = `Coin — учет личных финансов. Только чтение: расходы, доходы, наличность и транзакции пользователя.
Категории — это счета: expense (расходы), earnings (доходы); regular и debt — деньги в наличии и долги. Теги — подкатегории.
Расход — пополнение счета expense (amount_to), доход — списание со счета earnings (amount_from). Перевод между своими счетами — не доход и не расход.
Идентификаторы — UUID строками, их дают справочники get_accounts, get_tags, get_account_groups, get_currencies. Даты — yyyy-MM-dd, бакеты по UTC, неделя с понедельника.
Суммы — числа полной точности. Пересчет в валюту ответа идет по ТЕКУЩЕМУ курсу для всех дат: история курсов не хранится, поэтому не делай выводов по курсовым колебаниям.`

// NewServer создает MCP-сервер без инструментов: их добавляют модули
func NewServer(version string) *sdk.Server {
	server := sdk.NewServer(
		&sdk.Implementation{Name: "coin", Title: "Coin", Version: version}, //nolint:exhaustruct
		&sdk.ServerOptions{Instructions: instructions},                     //nolint:exhaustruct
	)
	server.AddReceivingMiddleware(logCalls)
	return server
}
