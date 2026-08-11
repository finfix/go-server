package accountBudgetDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaCoin + "." + "account_budgets"
	TableWithAlias = Table + " " + alias
	alias          = "ab"
)

const (
	ColumnID              = "id"
	ColumnAccountID       = "account_id"
	ColumnAccountGroupID  = "account_group_id"
	ColumnAmount          = "amount"
	ColumnFixedSum        = "fixed_sum"
	ColumnDaysOffset      = "days_offset"
	ColumnGradualFilling  = "gradual_filling"
	ColumnEffectiveFrom   = "effective_from"
	ColumnIsDeleted       = "is_deleted"
	ColumnCreatedByUserID = "created_by_user_id"
	ColumnDatetimeCreate  = "datetime_create"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
