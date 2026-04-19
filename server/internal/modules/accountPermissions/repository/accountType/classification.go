package accountType

type Classification int

const (
	AccountTypeByParent Classification = iota + 1
	AccountTypeByType
)

var ClassificationMatching = map[Type]Classification{
	Regular:   AccountTypeByType,
	Debt:      AccountTypeByType,
	Earnings:  AccountTypeByType,
	Expense:   AccountTypeByType,
	Balancing: AccountTypeByType,

	Parent:  AccountTypeByParent,
	General: AccountTypeByParent,
}
