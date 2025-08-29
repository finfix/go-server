package accountDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaCoin + "." + "accounts"
	TableWithAlias = Table + " " + alias
	alias          = "a"
)

const (
	ColumnID                   = "id"
	ColumnBudgetAmount         = "budget_amount"
	ColumnName                 = "name"
	ColumnIconID               = "icon_id"
	ColumnType                 = "type_signatura"
	ColumnCurrency             = "currency_signatura"
	ColumnVisible              = "visible"
	ColumnAccountGroupID       = "account_group_id"
	ColumnAccountingInHeader   = "accounting_in_header"
	ColumnParentAccountID      = "parent_account_id"
	ColumnSerialNumber         = "serial_number"
	ColumnBudgetGradualFilling = "budget_gradual_filling"
	ColumnIsParent             = "is_parent"
	ColumnBudgetFixedSum       = "budget_fixed_sum"
	ColumnBudgetDaysOffset     = "budget_days_offset"
	ColumnDatetimeCreate       = "datetime_create"
	ColumnCreatedByUserID      = "created_by_user_id"
	ColumnAccountingInCharts   = "accounting_in_charts"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
