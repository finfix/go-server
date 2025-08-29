package iconDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaCoin + "." + "icons"
	TableWithAlias = Table + " " + alias
	alias          = "i"
)

const (
	ColumnID   = "id"
	ColumnImg  = "img"
	ColumnName = "name"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
