package currencyDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaCoin + "." + "currencies"
	TableWithAlias = Table + " " + alias
	alias          = "c"
)

const (
	ColumnSlug   = "slug"
	ColumnName   = "name"
	ColumnRate   = "rate"
	ColumnSymbol = "symbol"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
