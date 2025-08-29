package transactionDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaCoin + "." + "transactions"
	TableWithAlias = Table + " " + alias
	alias          = "tr"
)

const (
	ColumnID                 = "id"
	ColumnDate               = "date_transaction"
	ColumnType               = "type_signatura"
	ColumnAmountFrom         = "amount_from"
	ColumnAmountTo           = "amount_to"
	ColumnNote               = "note"
	ColumnAccountFromID      = "account_from_id"
	ColumnAccountToID        = "account_to_id"
	ColumnIsExecuted         = "is_executed"
	ColumnDatetimeCreate     = "datetime_create"
	ColumnAccountingInCharts = "accounting_in_charts"
	ColumnCreatedByUserID    = "created_by_user_id"
	ColumnAccountGroupID     = "account_group_id"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
