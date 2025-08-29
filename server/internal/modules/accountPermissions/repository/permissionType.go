package repository

type permissionType int

const (
	ByType permissionType = iota + 1
	ByIsParent
)
