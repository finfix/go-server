package model

import "server/internal/modules/account/model/accountType"

type PermissionSet struct {
	TypeToPermissions     map[accountType.Type]AccountPermissions
	IsParentToPermissions map[bool]AccountPermissions
}
