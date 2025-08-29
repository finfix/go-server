package userDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaCoin + "." + "users"
	TableWithAlias = Table + " " + alias
	alias          = "u"
)

const (
	ColumnID              = "id"
	ColumnName            = "name"
	ColumnEmail           = "email"
	ColumnPasswordHash    = "password_hash"
	ColumnTimeCreate      = "time_create"
	ColumnDefaultCurrency = "default_currency_signatura"
	ColumnPasswordSalt    = "password_salt"
	ColumnIsAdmin         = "is_admin"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
