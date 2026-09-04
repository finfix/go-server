package pendingLinkedTransferDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaCoin + "." + "pending_linked_transfers"
	TableWithAlias = Table + " " + alias
	alias          = "plt"
)

const (
	ColumnID                   = "id"
	ColumnStatus               = "status"
	ColumnSourceTransactionID  = "source_transaction_id"
	ColumnSourceAccountID      = "source_account_id"
	ColumnTargetAccountID      = "target_account_id"
	ColumnAccountGroupID       = "account_group_id"
	ColumnDatetimeCreate       = "datetime_create"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
