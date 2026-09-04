package accountDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaCoin + "." + "accounts"
	TableWithAlias = Table + " " + alias
	alias          = "a"
)

const (
	ColumnID                 = "id"
	ColumnName               = "name"
	ColumnIconID             = "icon_id"
	ColumnType               = "account_type"
	ColumnCurrency           = "currency"
	ColumnVisible            = "visible"
	ColumnAccountGroupID     = "account_group_id"
	ColumnAccountingInHeader = "accounting_in_header"
	ColumnParentAccountID    = "parent_account_id"
	ColumnRank               = "rank"
	ColumnIsParent           = "is_parent"
	ColumnDatetimeCreate     = "datetime_create"
	ColumnCreatedByUserID    = "created_by_user_id"
	ColumnAccountingInCharts = "accounting_in_charts"
	ColumnIsDeleted          = "is_deleted"
	ColumnLinkedAccountID    = "linked_account_id"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
