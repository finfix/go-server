package model

import (
	"server/internal/enum/accountType"
)

type PermissionSet struct {
	TypeToPermissions     map[accountType.AccountType]AccountPermissions
	IsParentToPermissions map[bool]AccountPermissions
}
