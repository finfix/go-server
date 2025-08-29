package repository

import (
	"server/internal/ddl"
)

const (
	TableName          = ddl.SchemaPermissions + "." + "account_permissions"
	TableNameWithAlias = TableName + " " + alias
	alias              = "ap"
)

const (
	ColumnAccountType = "account_type"
	ColumnActionType  = "action_type"
	ColumnAccess      = "access"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
