package tagToTransactionDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaCoin + "." + "tags_to_transaction"
	TableWithAlias = Table + " " + alias
	alias          = "ttt"
)

const (
	ColumnTagID         = "tag_id"
	ColumnTransactionID = "transaction_id"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
