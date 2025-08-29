package versionDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaSettings + "." + "versions"
	TableWithAlias = Table + " " + alias
	alias          = "v"
)

const (
	ColumnName    = "name"
	ColumnVersion = "version"
)
