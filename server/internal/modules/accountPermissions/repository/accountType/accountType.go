package accountType

import "server/internal/modules/account/model/accountType"

type Type string

const (
	Parent  Type = "parent"
	General Type = "general"

	Regular   = Type(accountType.Regular)
	Debt      = Type(accountType.Debt)
	Earnings  = Type(accountType.Earnings)
	Expense   = Type(accountType.Expense)
	Balancing = Type(accountType.Balancing)
)
