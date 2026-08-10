package accountGroupDDL

import "server/internal/ddl"

const (
	TableName          = ddl.SchemaCoin + "." + "account_groups"
	TableNameWithAlias = TableName + " " + alias
	alias              = "ag"
)

const (
	ColumnID              = "id"
	ColumnName            = "name"
	ColumnAvailableBudget = "available_budget"
	ColumnCurrency        = "currency_signatura"
	ColumnSerialNumber    = "serial_number"
	ColumnVisible         = "visible"
	ColumnDatetimeCreate  = "datetime_create"
	ColumnCreatedByUserID = "created_by_user_id"
	ColumnIsDeleted       = "is_deleted"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
