package tagDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaCoin + "." + "tags"
	TableWithAlias = Table + " " + alias
	alias          = "t"
)

const (
	ColumnID              = "id"
	ColumnName            = "name"
	ColumnAccountGroupID  = "account_group_id"
	ColumnCreatedByUserID = "created_by_user_id"
	ColumnDatetimeCreate  = "datetime_create"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
