package currencyRateDDL

import "server/internal/ddl"

const (
	Table = ddl.SchemaCoin + "." + "currency_rates"
)

const (
	ColumnCurrency       = "currency"         // Код валюты
	ColumnValidAt        = "valid_at"         // С какого момента действует курс
	ColumnRate           = "rate"             // Единиц валюты за 1 USD
	ColumnSource         = "source"           // provider или manual
	ColumnAccountGroupID = "account_group_id" // Группа; NULL — курс для всех
)

// Источники курса
const (
	SourceProvider = "provider" // Провайдер курсов
	SourceManual   = "manual"   // Введен вручную
)
