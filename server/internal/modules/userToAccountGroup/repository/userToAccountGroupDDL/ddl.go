package userToAccountGroupDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaCoin + "." + "users_to_account_groups"
	TableWithAlias = Table + " " + alias
	alias          = "utag"
)

const (
	ColumnUserID         = "user_id"
	ColumnAccountGroupID = "account_group_id"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
